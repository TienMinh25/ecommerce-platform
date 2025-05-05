import { ChevronDownIcon, HamburgerIcon, SearchIcon } from '@chakra-ui/icons';
import {
  Avatar,
  Badge,
  Box,
  Button,
  Container,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Flex,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Text,
  useBreakpointValue,
  useDisclosure,
  VStack,
} from '@chakra-ui/react';
import { useState } from 'react';
import { FaShoppingCart, FaStore, FaShippingFast, FaUserShield } from 'react-icons/fa';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import Logo from '../ui/Logo';
import NotificationBell from '../notifications/NotificationBell';

// Define role constants
const ROLE_CUSTOMER = 'customer';
const ROLE_SUPPLIER = 'supplier';
const ROLE_DELIVERER = 'deliverer';
const ROLE_ADMIN = 'admin';

const Header = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();
  const { user, logout } = useAuth();

  const isMobile = useBreakpointValue({ base: true, md: false });

  // Handle search in Header component
  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      // Navigate to product listing page with search query
      navigate({
        pathname: '/products',
        search: `?keyword=${encodeURIComponent(searchQuery)}`,
      });

      // Reset search query after navigation
      setSearchQuery('');

      // Close mobile drawer if open
      if (isOpen) {
        onClose();
      }
    }
  };

  // Handle mobile search
  const handleMobileSearch = () => {
    // Implement mobile search functionality
    if (searchQuery.trim()) {
      navigate(`/products?keyword=${encodeURIComponent(searchQuery)}`);
      setSearchQuery('');
      onClose();
    }
  };

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login', {replace: true});
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  // Check if user has specific roles
  const isSupplier = user?.hasRole(ROLE_SUPPLIER) || false;
  const isDeliverer = user?.hasRole(ROLE_DELIVERER) || false;
  const isAdmin = user?.hasRole(ROLE_ADMIN) || false;

  return (
      <Box
          as='header'
          py={2}
          borderBottom='1px'
          borderColor='gray.200'
          bg='white'
          boxShadow='sm'
          position='sticky'
          top='0'
          zIndex='999'
      >
        <Container maxW='container.xl'>
          <Flex align='center' justify='space-between'>
            {/* Logo and Mobile Menu Button */}
            <HStack spacing={4}>
              {isMobile && (
                  <IconButton
                      icon={<HamburgerIcon />}
                      onClick={onOpen}
                      variant='ghost'
                      aria-label='Menu'
                  />
              )}
              <Logo size={isMobile ? 'sm' : 'md'} />
            </HStack>

            {/* Search */}
            {!isMobile && (
                <Box flex='1' mx={8}>
                  <form onSubmit={handleSearch}>
                    <InputGroup size='md'>
                      <Input
                          placeholder='Tìm kiếm sản phẩm...'
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                          bg='gray.50'
                          _focus={{
                            borderColor: 'brand.500',
                            boxShadow: '0 0 0 1px var(--chakra-colors-brand-500)',
                          }}
                      />
                      <InputRightElement>
                        <IconButton
                            aria-label='Search'
                            icon={<SearchIcon />}
                            type='submit'
                            variant='ghost'
                            colorScheme='brand'
                        />
                      </InputRightElement>
                    </InputGroup>
                  </form>
                </Box>
            )}

            {/* Nav Icons */}
            <HStack spacing={{ base: 2, md: 4 }}>
              {isMobile && (
                  <IconButton
                      aria-label='Search'
                      icon={<SearchIcon />}
                      onClick={handleMobileSearch}
                      variant='ghost'
                  />
              )}

              {/* Notification bell */}
              <Box display={{ base: 'none', md: 'block' }}>
                <NotificationBell />
              </Box>

              {/* Cart icon */}
              <Box position='relative'>
                <IconButton
                    as={RouterLink}
                    to='/cart'
                    aria-label='Shopping Cart'
                    icon={<FaShoppingCart />}
                    variant='ghost'
                />
                <Badge
                    position='absolute'
                    top='-2px'
                    right='-2px'
                    colorScheme='brand'
                    borderRadius='full'
                    size='xs'
                >
                  2
                </Badge>
              </Box>

              {/* User menu */}
              <Menu>
                <MenuButton
                    as={Button}
                    rightIcon={<ChevronDownIcon />}
                    variant='ghost'
                    display={{ base: 'none', md: 'flex' }}
                    minW="auto"
                    h="40px"
                    px={3}
                    _hover={{
                      bg: 'gray.100',
                    }}
                    _active={{
                      bg: 'gray.200',
                    }}
                >
                  <Flex align="center" h="100%">
                    <Avatar
                        size='sm'
                        name={user?.fullname || 'User'}
                        src={user?.avatarUrl}
                        mr={2}
                    />
                    <Text
                        display={{ base: 'none', lg: 'block' }}
                        fontWeight="medium"
                        fontSize="sm"
                        mt="1px"
                    >
                      {user?.fullname || 'Tài khoản'}
                    </Text>
                  </Flex>
                </MenuButton>
                <MenuList
                    zIndex={1000}
                    p={0}
                    overflow="hidden"
                    borderRadius="md"
                    boxShadow="lg"
                >
                  <MenuItem
                      as={RouterLink}
                      to='/user/account/profile'
                      py={3}
                      _hover={{ bg: 'gray.50' }}
                      w="full"
                  >
                    Tài khoản của tôi
                  </MenuItem>
                  <MenuItem
                      as={RouterLink}
                      to='/account/orders'
                      py={3}
                      _hover={{ bg: 'gray.50' }}
                      w="full"
                  >
                    Đơn hàng của tôi
                  </MenuItem>

                  {isSupplier && (
                      <>
                        <MenuDivider m={0} />
                        <MenuItem
                            as={RouterLink}
                            to='/supplier/dashboard'
                            py={3}
                            _hover={{ bg: 'gray.50' }}
                            w="full"
                        >
                          <Flex align="center">
                            <FaStore style={{ marginRight: '8px' }} />
                            Quản lý cửa hàng
                          </Flex>
                        </MenuItem>
                      </>
                  )}

                  {isDeliverer && (
                      <>
                        <MenuDivider m={0} />
                        <MenuItem
                            as={RouterLink}
                            to='/delivery/orders'
                            py={3}
                            _hover={{ bg: 'gray.50' }}
                            w="full"
                        >
                          <Flex align="center">
                            <FaShippingFast style={{ marginRight: '8px' }} />
                            Quản lý giao hàng
                          </Flex>
                        </MenuItem>
                      </>
                  )}

                  {isAdmin && (
                      <>
                        <MenuDivider m={0} />
                        <MenuItem
                            as={RouterLink}
                            to='/dashboard'
                            py={3}
                            _hover={{ bg: 'gray.50' }}
                            w="full"
                        >
                          <Flex align="center">
                            <FaUserShield style={{ marginRight: '8px' }} />
                            Trang quản trị
                          </Flex>
                        </MenuItem>
                      </>
                  )}

                  <MenuDivider m={0} />
                  <MenuItem
                      onClick={handleLogout}
                      color="red.500"
                      py={3}
                      _hover={{ bg: 'red.50' }}
                      w="full"
                  >
                    Đăng xuất
                  </MenuItem>
                </MenuList>
              </Menu>
            </HStack>
          </Flex>
        </Container>

        {/* Mobile Drawer Menu */}
        <Drawer isOpen={isOpen} placement='left' onClose={onClose}>
          <DrawerOverlay />
          <DrawerContent>
            <DrawerCloseButton />
            <DrawerHeader borderBottomWidth="1px" py={3} px={4}>
              <Logo />
            </DrawerHeader>

            <DrawerBody px={0}>
              <VStack align='stretch' spacing={0}>
                {/* User profile in drawer */}
                <Flex
                    alignItems='center'
                    p={4}
                    bg="gray.50"
                    borderBottomWidth="1px"
                >
                  <Avatar
                      size='md'
                      name={user?.fullname}
                      src={user?.avatarUrl}
                      mr={3}
                  />
                  <Box>
                    <Text fontWeight='bold'>{user?.fullname}</Text>
                    <Text fontSize='sm' color='gray.600'>
                      {user?.email || user?.fullname}
                    </Text>
                  </Box>
                </Flex>

                <Box px={4} py={3}>
                  <Box as='form' onSubmit={handleSearch}>
                    <InputGroup>
                      <Input
                          placeholder='Tìm kiếm sản phẩm...'
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                      />
                      <InputRightElement>
                        <IconButton
                            aria-label='Search'
                            icon={<SearchIcon />}
                            type='submit'
                            variant='ghost'
                        />
                      </InputRightElement>
                    </InputGroup>
                  </Box>
                </Box>

                {/* Product categories */}
                <Box py={2} px={4} bg="gray.100" fontWeight="bold" color="gray.700">
                  Danh mục
                </Box>

                <Box
                    as={RouterLink}
                    to='/'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Trang chủ
                </Box>

                <Box
                    as={RouterLink}
                    to='/products'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Sản phẩm
                </Box>

                <Box
                    as={RouterLink}
                    to='/categories'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Danh mục
                </Box>

                <Box
                    as={RouterLink}
                    to='/promotions'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Khuyến mãi
                </Box>

                {/* Account */}
                <Box py={2} px={4} bg="gray.100" fontWeight="bold" color="gray.700">
                  Tài khoản
                </Box>

                <Box
                    as={RouterLink}
                    to='/account/profile'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Hồ sơ cá nhân
                </Box>

                <Box
                    as={RouterLink}
                    to='/account/orders'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Đơn hàng của tôi
                </Box>

                <Box
                    as={RouterLink}
                    to='/cart'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Giỏ hàng
                </Box>

                <Box
                    as={RouterLink}
                    to='/notifications'
                    py={3}
                    px={4}
                    borderBottomWidth="1px"
                    bg="white"
                    _hover={{ bg: 'gray.50' }}
                    onClick={onClose}
                    fontWeight="medium"
                    w="full"
                >
                  Thông báo
                </Box>

                {isSupplier && (
                    <Box
                        as={RouterLink}
                        to='/supplier/dashboard'
                        py={3}
                        px={4}
                        borderBottomWidth="1px"
                        bg="white"
                        _hover={{ bg: 'gray.50' }}
                        onClick={onClose}
                        fontWeight="medium"
                        w="full"
                    >
                      <Flex align="center">
                        <FaStore style={{ marginRight: '8px' }} />
                        Quản lý cửa hàng
                      </Flex>
                    </Box>
                )}

                {isDeliverer && (
                    <Box
                        as={RouterLink}
                        to='/delivery/orders'
                        py={3}
                        px={4}
                        borderBottomWidth="1px"
                        bg="white"
                        _hover={{ bg: 'gray.50' }}
                        onClick={onClose}
                        fontWeight="medium"
                        w="full"
                    >
                      <Flex align="center">
                        <FaShippingFast style={{ marginRight: '8px' }} />
                        Quản lý giao hàng
                      </Flex>
                    </Box>
                )}

                {isAdmin && (
                    <Box
                        as={RouterLink}
                        to='/admin/dashboard'
                        py={3}
                        px={4}
                        borderBottomWidth="1px"
                        bg="white"
                        _hover={{ bg: 'gray.50' }}
                        onClick={onClose}
                        fontWeight="medium"
                        w="full"
                    >
                      <Flex align="center">
                        <FaUserShield style={{ marginRight: '8px' }} />
                        Trang quản trị
                      </Flex>
                    </Box>
                )}

                <Box
                    as='button'
                    py={3}
                    px={4}
                    color='red.500'
                    textAlign='left'
                    w="full"
                    bg="white"
                    _hover={{ bg: 'red.50' }}
                    fontWeight="medium"
                    onClick={() => {
                      handleLogout();
                      onClose();
                    }}
                >
                  Đăng xuất
                </Box>
              </VStack>
            </DrawerBody>
          </DrawerContent>
        </Drawer>
      </Box>
  );
};

export default Header;