import { useState, useEffect, useRef } from 'react';
import {
  Box,
  Flex,
  Image,
  IconButton,
  Heading,
  useBreakpointValue,
} from '@chakra-ui/react';
import { FaChevronLeft, FaChevronRight } from 'react-icons/fa';

const PaymentLogosCarousel = ({ logos }) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isPaused, setIsPaused] = useState(false);
  const logoContainerRef = useRef(null);
  const showCount = useBreakpointValue({ base: 2, sm: 3, md: 4 });
  const maxIndex = logos.length - showCount;

  // Tự động chuyển slide mỗi 3 giây
  useEffect(() => {
    if (isPaused) return;

    const interval = setInterval(() => {
      setCurrentIndex((prevIndex) =>
        prevIndex >= maxIndex ? 0 : prevIndex + 1,
      );
    }, 3000);

    return () => clearInterval(interval);
  }, [maxIndex, isPaused]);

  const handlePrev = () => {
    setCurrentIndex((prevIndex) => (prevIndex <= 0 ? maxIndex : prevIndex - 1));
  };

  const handleNext = () => {
    setCurrentIndex((prevIndex) => (prevIndex >= maxIndex ? 0 : prevIndex + 1));
  };

  return (
    <Box position='relative' py={4}>
      <Heading as='h5' size='sm' mb={4} textAlign='center'>
        Chúng tôi chấp nhận thanh toán qua
      </Heading>

      <Flex alignItems='center'>
        <IconButton
          icon={<FaChevronLeft />}
          onClick={handlePrev}
          aria-label='Previous'
          variant='ghost'
          colorScheme='gray'
          size='sm'
        />

        <Box
          overflow='hidden'
          mx={2}
          flex='1'
          onMouseEnter={() => setIsPaused(true)}
          onMouseLeave={() => setIsPaused(false)}
        >
          <Flex
            ref={logoContainerRef}
            transition='transform 0.5s ease'
            transform={`translateX(-${currentIndex * (100 / showCount)}%)`}
          >
            {logos.map((logo, index) => (
              <Box
                key={index}
                flexBasis={`${100 / showCount}%`}
                flexShrink={0}
                p={2}
              >
                <Box
                  borderWidth='1px'
                  borderRadius='md'
                  boxShadow='sm'
                  p={3}
                  bg='white'
                  height='60px'
                  display='flex'
                  alignItems='center'
                  justifyContent='center'
                  transition='all 0.3s'
                  _hover={{
                    transform: 'translateY(-2px)',
                    boxShadow: 'md',
                    borderColor: 'gray.300',
                  }}
                >
                  <Image
                    src={logo.src}
                    alt={logo.alt}
                    maxH='40px'
                    maxW='100%'
                    objectFit='contain'
                  />
                </Box>
              </Box>
            ))}
          </Flex>
        </Box>

        <IconButton
          icon={<FaChevronRight />}
          onClick={handleNext}
          aria-label='Next'
          variant='ghost'
          colorScheme='gray'
          size='sm'
        />
      </Flex>

      {/* Indicators */}
      <Flex justify='center' mt={3}>
        {Array.from({ length: maxIndex + 1 }).map((_, index) => (
          <Box
            key={index}
            w='8px'
            h='8px'
            borderRadius='full'
            mx={1}
            bg={currentIndex === index ? 'brand.500' : 'gray.300'}
            cursor='pointer'
            onClick={() => setCurrentIndex(index)}
          />
        ))}
      </Flex>
    </Box>
  );
};

export default PaymentLogosCarousel;
