import React, { useRef, useState, useEffect } from 'react';
import { Box, Flex, Text, Image, IconButton, useBreakpointValue, Heading } from '@chakra-ui/react';
import { ChevronLeftIcon, ChevronRightIcon } from '@chakra-ui/icons';
import categoryService from '../../services/categoryService';

const CategorySlider = () => {
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const sliderRef = useRef(null);
    const [scrollPosition, setScrollPosition] = useState(0);
    const [maxScroll, setMaxScroll] = useState(0);

    // Fetch categories from API
    useEffect(() => {
        const fetchCategories = async () => {
            try {
                setLoading(true);
                // Sử dụng service để gọi API
                const response = await categoryService.getAllCategories();
                if (response.data && response.data.data) {
                    setCategories(response.data.data);
                }
            } catch (err) {
                console.error('Error fetching categories:', err);
                setError('Không thể tải danh mục sản phẩm');
            } finally {
                setLoading(false);
            }
        };

        fetchCategories();
    }, []);

    // Calculate max scroll when categories or window size changes
    useEffect(() => {
        const calculateMaxScroll = () => {
            if (sliderRef.current) {
                const containerWidth = sliderRef.current.clientWidth;
                const scrollWidth = sliderRef.current.scrollWidth;
                setMaxScroll(Math.max(0, scrollWidth - containerWidth));
            }
        };

        calculateMaxScroll();
        window.addEventListener('resize', calculateMaxScroll);

        return () => window.removeEventListener('resize', calculateMaxScroll);
    }, [categories]);

    const scrollLeft = () => {
        if (sliderRef.current) {
            const newPosition = Math.max(0, scrollPosition - 300);
            sliderRef.current.scrollTo({ left: newPosition, behavior: 'smooth' });
            setScrollPosition(newPosition);
        }
    };

    const scrollRight = () => {
        if (sliderRef.current) {
            const newPosition = Math.min(maxScroll, scrollPosition + 300);
            sliderRef.current.scrollTo({ left: newPosition, behavior: 'smooth' });
            setScrollPosition(newPosition);
        }
    };

    // Handle manual scroll
    const handleScroll = () => {
        if (sliderRef.current) {
            setScrollPosition(sliderRef.current.scrollLeft);
        }
    };

    const showControls = useBreakpointValue({ base: false, md: true });
    const itemWidth = useBreakpointValue({ base: '150px', md: '180px', lg: '200px' });

    if (loading) {
        return (
            <Box py={6} textAlign="center">
                Đang tải danh mục...
            </Box>
        );
    }

    if (error) {
        return (
            <Box py={6} textAlign="center" color="red.500">
                {error}
            </Box>
        );
    }

    // Placeholder categories nếu API chưa hoạt động hoặc không có dữ liệu
    const placeholderCategories = [
        {
            id: 1,
            name: 'Điện tử & Công nghệ',
            image_url: 'https://images.unsplash.com/photo-1519389950473-47ba0277781c?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80'
        },
        {
            id: 2,
            name: 'Thời trang',
            image_url: 'https://images.unsplash.com/photo-1551232864-3f0890e580d9?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1974&q=80'
        },
        {
            id: 3,
            name: 'Nhà cửa & Đời sống',
            image_url: 'https://images.unsplash.com/photo-1583847268964-b28dc8f51f92?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1974&q=80'
        },
        {
            id: 4,
            name: 'Sách & Văn phòng phẩm',
            image_url: 'https://images.unsplash.com/photo-1524578271613-d550eacf6090?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80'
        },
        {
            id: 5,
            name: 'Thể thao & Du lịch',
            image_url: 'https://images.unsplash.com/photo-1530549387789-4c1017266635?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=2070&q=80'
        }
    ];

    // Sử dụng dữ liệu từ API hoặc placeholder nếu không có dữ liệu
    const displayCategories = categories.length > 0 ? categories : placeholderCategories;

    return (
        <Box position="relative" my={8} py={4}>
            <Heading as="h2" size="lg" mb={6} px={4} fontWeight="700" color="gray.800">
                Danh mục nổi bật
            </Heading>

            <Box position="relative">
                {showControls && scrollPosition > 0 && (
                    <IconButton
                        aria-label="Scroll left"
                        icon={<ChevronLeftIcon boxSize={6} color={'black'}/>}
                        onClick={scrollLeft}
                        position="absolute"
                        left={2}
                        top="50%"
                        transform="translateY(-50%)"
                        zIndex={2}
                        borderRadius="full"
                        bg="white"
                        boxShadow="md"
                        _hover={{ bg: 'gray.100' }}
                    />
                )}

                <Box
                    ref={sliderRef}
                    overflowX="auto"
                    css={{
                        '&::-webkit-scrollbar': {
                            display: 'none',
                        },
                        scrollbarWidth: 'none',
                    }}
                    onScroll={handleScroll}
                    px={4}
                >
                    <Flex gap={4}>
                        {displayCategories.map((category) => (
                            <Box
                                key={category.id}
                                as="a"
                                href={`/products?categoryID=${category.id}`}
                                minW={itemWidth}
                                maxW={itemWidth}
                                textAlign="center"
                                transition="all 0.3s"
                                _hover={{ transform: 'translateY(-5px)' }}
                                borderRadius="lg"
                                overflow="hidden"
                                boxShadow="sm"
                                height="auto"
                            >
                                <Box
                                    overflow="hidden"
                                    position="relative"
                                    height="150px"
                                >
                                    <Image
                                        src={category.image_url || "https://via.placeholder.com/200x150"}
                                        alt={category.name}
                                        width="100%"
                                        height="100%"
                                        objectFit="cover"
                                        transition="transform 0.3s ease"
                                        _hover={{ transform: 'scale(1.05)' }}
                                    />
                                </Box>

                                {/* Phần tên danh mục với nền xanh */}
                                <Box
                                    p={3}
                                    bg="blue.500"
                                    color="white"
                                    textAlign="center"
                                >
                                    <Text
                                        fontWeight="600"
                                        fontSize="md"
                                        isTruncated
                                    >
                                        {category.name}
                                    </Text>
                                </Box>
                            </Box>
                        ))}
                    </Flex>
                </Box>

                {showControls && scrollPosition < maxScroll && (
                    <IconButton
                        aria-label="Scroll right"
                        icon={<ChevronRightIcon boxSize={6} color={'black'}/>}
                        onClick={scrollRight}
                        position="absolute"
                        right={2}
                        top="50%"
                        transform="translateY(-50%)"
                        zIndex={2}
                        borderRadius="full"
                        bg="white"
                        boxShadow="md"
                        _hover={{ bg: 'gray.100' }}
                    />
                )}
            </Box>
        </Box>
    );
};

export default CategorySlider;