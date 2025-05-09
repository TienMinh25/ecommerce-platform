import React from 'react';
import {
    Box,
    Flex,
    Heading,
    Text,
    Avatar,
    HStack,
    VStack,
    Skeleton,
    SkeletonCircle,
    Icon,
    Center,
} from '@chakra-ui/react';
import { StarIcon } from '@chakra-ui/icons';
import { FaThumbsUp } from 'react-icons/fa';
import Pagination from '../common/Pagination';

const RatingSummary = ({ productRating, totalReviews }) => {
    return (
        <Flex
            direction={{ base: 'column', md: 'row' }}
            justify="space-between"
            align={{ base: 'center', md: 'flex-start' }}
            p={4}
            borderWidth="1px"
            borderRadius="md"
            bg="gray.50"
            mb={6}
        >
            <Box textAlign={{ base: 'center', md: 'left' }} mb={{ base: 4, md: 0 }}>
                <Text fontSize="3xl" fontWeight="bold">
                    {productRating.toFixed(1)}/5
                </Text>
                <HStack spacing={1} justify={{ base: 'center', md: 'flex-start' }}>
                    {Array(5)
                        .fill('')
                        .map((_, i) => (
                            <StarIcon
                                key={i}
                                color={i < Math.round(productRating) ? 'yellow.400' : 'gray.300'}
                            />
                        ))}
                </HStack>
                <Text color="gray.600" fontSize="sm">
                    {totalReviews} đánh giá
                </Text>
            </Box>
        </Flex>
    );
};

const ReviewItem = ({ review }) => {
    // Format date to DD/MM/YYYY
    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return `${date.getDate().toString().padStart(2, '0')}/${(date.getMonth() + 1).toString().padStart(2, '0')}/${date.getFullYear()}`;
    };

    return (
        <Box p={4} borderWidth="1px" borderRadius="md" mb={4}>
            <Flex mb={4}>
                <Avatar
                    size="sm"
                    name={review.user_name}
                    src={review.user_avatar_url}
                    mr={4}
                />
                <Box flex="1">
                    <Flex justify="space-between" align="center">
                        <Text fontWeight="bold">{review.user_name}</Text>
                        <Text fontSize="sm" color="gray.500">
                            {formatDate(review.created_at)}
                        </Text>
                    </Flex>
                    <HStack spacing={1}>
                        {Array(5)
                            .fill('')
                            .map((_, i) => (
                                <StarIcon
                                    key={i}
                                    size="sm"
                                    color={i < Math.round(review.rating) ? 'yellow.400' : 'gray.300'}
                                />
                            ))}
                        <Text fontSize="sm" color="gray.500" ml={1}>
                            {review.rating.toFixed(1)}
                        </Text>
                    </HStack>
                </Box>
            </Flex>
            <Text>{review.comment}</Text>
            {review.helpful_votes > 0 && (
                <Flex align="center" mt={2}>
                    <Icon as={FaThumbsUp} color="gray.500" mr={1} />
                    <Text fontSize="sm" color="gray.500">
                        {review.helpful_votes} người thấy hữu ích
                    </Text>
                </Flex>
            )}
        </Box>
    );
};

const ReviewSkeleton = () => (
    <Box p={4} borderWidth="1px" borderRadius="md" mb={4}>
        <Flex mb={4}>
            <SkeletonCircle size="8" mr={4} />
            <Box flex="1">
                <Skeleton height="20px" width="120px" mb={2} />
                <Skeleton height="16px" width="100px" />
            </Box>
            <Skeleton height="16px" width="80px" />
        </Flex>
        <Skeleton height="20px" mb={2} />
        <Skeleton height="20px" width="90%" />
    </Box>
);

const ProductReviews = ({
                            reviews,
                            metadata,
                            isLoading,
                            onPageChange,
                            productRating,
                            totalReviews,
                        }) => {
    return (
        <Box>
            <Heading as="h3" size="md" mb={6}>
                Đánh giá từ khách hàng
            </Heading>

            {/* Rating Summary */}
            <RatingSummary productRating={productRating} totalReviews={totalReviews} />

            {/* Reviews List */}
            {isLoading ? (
                <VStack spacing={4} align="stretch">
                    {[1, 2, 3].map((_, index) => (
                        <ReviewSkeleton key={index} />
                    ))}
                </VStack>
            ) : reviews.length > 0 ? (
                <VStack spacing={4} align="stretch">
                    {reviews.map((review) => (
                        <ReviewItem key={review.id} review={review} />
                    ))}

                    {/* Pagination */}
                    {metadata && metadata.total_pages > 1 && (
                        <Box mt={8} display="flex" justifyContent="center">
                            <Pagination
                                currentPage={metadata.page}
                                totalPages={metadata.total_pages}
                                onPageChange={onPageChange}
                            />
                        </Box>
                    )}
                </VStack>
            ) : (
                <Center py={10}>
                    <VStack>
                        <Text color="gray.500">Chưa có đánh giá nào cho sản phẩm này</Text>
                        <Text color="gray.400" fontSize="sm">Đánh giá chỉ có thể viết sau khi mua hàng</Text>
                    </VStack>
                </Center>
            )}
        </Box>
    );
};

export default ProductReviews;