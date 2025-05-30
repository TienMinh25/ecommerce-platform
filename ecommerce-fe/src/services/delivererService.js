import api from './api';

const delivererService = {
    /**
     * Get presigned URL for S3 upload
     * @param {Object} fileData - File information
     * @returns {Promise} Promise object with presigned URL
     */
    getPresignedUrl: async (fileData) => {
        try {
            const response = await api.post('/s3/presigned_url', {
                file_name: fileData.fileName,
                content_type: fileData.contentType,
                file_size: fileData.fileSize,
                bucket_name: fileData.bucketName || 'deliverers'
            });
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Upload file to S3 using presigned URL
     * @param {string} presignedUrl - Presigned URL from S3
     * @param {File} file - File to upload
     * @returns {Promise} Promise object with upload result
     */
    uploadToS3: async (presignedUrl, file) => {
        try {
            const response = await fetch(presignedUrl, {
                method: 'PUT',
                body: file,
                headers: {
                    'Content-Type': file.type,
                },
            });

            if (!response.ok) {
                throw new Error(`Upload failed: ${response.statusText}`);
            }

            // Return the S3 URL (remove query parameters)
            return presignedUrl.split('?')[0];
        } catch (error) {
            throw error;
        }
    },

    /**
     * Complete file upload process (get presigned URL + upload)
     * @param {File} file - File to upload
     * @param {string} bucketName - S3 bucket name
     * @returns {Promise} Promise object with S3 URL
     */
    uploadFile: async (file, bucketName = 'deliverers') => {
        try {
            // Get presigned URL
            const { url } = await delivererService.getPresignedUrl({
                fileName: file.name,
                contentType: file.type,
                fileSize: file.size,
                bucketName: bucketName
            });

            // Upload file to S3
            const s3Url = await delivererService.uploadToS3(url, file);

            return s3Url;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Register as deliverer
     * @param {Object} delivererData - Deliverer registration data
     * @returns {Promise} Promise object with registration result
     */
    registerDeliverer: async (delivererData) => {
        try {
            const response = await api.post('/deliverers/register', delivererData);
            return response.data.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Upload multiple files and return their URLs
     * @param {Object} files - Object containing files to upload
     * @param {string} bucketName - S3 bucket name
     * @returns {Promise} Promise object with URLs mapping
     */
    uploadMultipleFiles: async (files, bucketName = 'deliverers') => {
        try {
            const uploadPromises = [];
            const uploadedUrls = {};

            for (const [key, file] of Object.entries(files)) {
                if (file) {
                    uploadPromises.push(
                        delivererService.uploadFile(file, bucketName).then(url => {
                            uploadedUrls[key] = url;
                        })
                    );
                }
            }

            await Promise.all(uploadPromises);
            return uploadedUrls;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Complete deliverer registration process (upload files + register)
     * @param {Object} formData - Form data
     * @param {Object} files - Files to upload
     * @returns {Promise} Promise object with registration result
     */
    completeRegistration: async (formData, files) => {
        try {
            // Upload driving license images first
            const uploadedUrls = await delivererService.uploadMultipleFiles({
                drivingLicenseFront: files.drivingLicenseFront,
                drivingLicenseBack: files.drivingLicenseBack
            }, 'deliverers');

            // Prepare submission data according to API schema
            const submissionData = {
                id_card_number: formData.drivingLicenseNumber,
                vehicle_type: formData.vehicleType,
                vehicle_license_plate: formData.vehicleLicensePlate,
                id_card_front_image: uploadedUrls.drivingLicenseFront,
                id_card_back_image: uploadedUrls.drivingLicenseBack,
                service_area: {
                    country: 'Việt Nam',
                    city: formData.selectedProvinceName,
                    district: formData.selectedDistrictName,
                    ward: formData.selectedWardName
                }
            };

            // Register deliverer
            const result = await delivererService.registerDeliverer(submissionData);
            return { success: true, data: result };
        } catch (error) {
            return { success: false, error: error };
        }
    }
};

export default delivererService;