import api from './api';

const supplierService = {
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
                bucket_name: fileData.bucketName || 'suppliers'
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
    uploadFile: async (file, bucketName = 'suppliers') => {
        try {
            // Get presigned URL
            const { url } = await supplierService.getPresignedUrl({
                fileName: file.name,
                contentType: file.type,
                fileSize: file.size,
                bucketName: bucketName
            });

            // Upload file to S3
            const s3Url = await supplierService.uploadToS3(url, file);

            return s3Url;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Register as supplier
     * @param {Object} supplierData - Supplier registration data
     * @returns {Promise} Promise object with registration result
     */
    registerSupplier: async (supplierData) => {
        try {
            const response = await api.post('/suppliers/register', supplierData);
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
    uploadMultipleFiles: async (files, bucketName = 'suppliers') => {
        try {
            const uploadPromises = [];
            const uploadedUrls = {};

            for (const [key, file] of Object.entries(files)) {
                if (file) {
                    uploadPromises.push(
                        supplierService.uploadFile(file, bucketName).then(url => {
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
     * Complete supplier registration process (upload files + register)
     * @param {Object} formData - Form data
     * @param {Object} files - Files to upload
     * @param {Object} selectedAddress - Selected business address
     * @returns {Promise} Promise object with registration result
     */
    completeRegistration: async (formData, files, selectedAddress) => {
        try {
            // Upload all files first
            const uploadedUrls = await supplierService.uploadMultipleFiles(files, 'suppliers');

            // Prepare submission data according to API schema
            const submissionData = {
                company_name: formData.companyName,
                contact_phone: formData.contactPhone,
                tax_id: formData.taxId,
                description: formData.description,
                logo_company_url: uploadedUrls.logo,
                business_address_id: selectedAddress.id,
                documents: {
                    business_license: uploadedUrls.business_license,
                    tax_certificate: uploadedUrls.tax_certificate,
                    id_card_front: uploadedUrls.id_card_front,
                    id_card_back: uploadedUrls.id_card_back
                }
            };

            // Register supplier
            const result = await supplierService.registerSupplier(submissionData);
            return { success: true, data: result };
        } catch (error) {
            return { success: false, error: error };
        }
    },

    /**
     * Get suppliers with filters and pagination
     * @param {Object} params - Query parameters
     * @returns {Promise} Promise object with suppliers data
     */
    getSuppliers: async (params) => {
        try {
            const response = await api.get('/suppliers', { params });
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Get supplier by ID
     * @param {string} supplierId - Supplier ID
     * @returns {Promise} Promise object with supplier details
     */
    getSupplierById: async (supplierId) => {
        try {
            const response = await api.get(`/suppliers/${supplierId}`);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Update supplier status
     * @param {string} supplierId - Supplier ID
     * @param {Object} updateData - Update data
     * @returns {Promise} Promise object with update result
     */
    updateSupplier: async (supplierId, updateData) => {
        try {
            const response = await api.patch(`/suppliers/${supplierId}`, updateData);
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Update supplier document verification status (approve/reject)
     * @param {string} supplierId - Supplier ID
     * @param {string} documentId - Document ID
     * @param {string} status - Document status ('approved', 'rejected', 'pending')
     * @returns {Promise} Promise object with update result
     */
    updateDocumentVerificationStatus: async (supplierId, documentId, status) => {
        try {
            const response = await api.patch(`/suppliers/${supplierId}/documents/${documentId}`, {
                status: status
            });
            return response.data;
        } catch (error) {
            throw error;
        }
    },

    /**
     * Get supplier orders (Supplier's own orders)
     * @param {Object} params - Query parameters
     * @param {number} params.page - Page number
     * @param {number} params.limit - Items per page
     * @param {string} params.status - Order status filter (based on acceptStatus)
     * @param {string} params.keyword - Search keyword
     * @returns {Promise} Promise object with supplier orders data
     */
    getSupplierOrders: async (params) => {
        try {
            const response = await api.get('/suppliers/me', { params });
            return response.data;
        } catch (error) {
            throw error;
        }
    },
};

export default supplierService;