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
import { FaBell, FaHeart, FaShoppingCart } from 'react-icons/fa';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import Logo from '../ui/Logo';

const Header = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();
  const { isAuthenticated, user, logout } = useAuth();

  const isMobile = useBreakpointValue({ base: true, md: false });

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      // TODO: Implement search
      navigate(`/products?search=${encodeURIComponent(searchQuery)}`);
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

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
                onClick={() => {
                  /* TODO: Open mobile search */
                }}
                variant='ghost'
              />
            )}

            {isAuthenticated ? (
              <>
                <IconButton
                  as={RouterLink}
                  to='/favorites'
                  aria-label='Favorites'
                  icon={<FaHeart />}
                  variant='ghost'
                  display={{ base: 'none', md: 'flex' }}
                />

                <Box
                  position='relative'
                  display={{ base: 'none', md: 'block' }}
                >
                  <IconButton
                    as={RouterLink}
                    to='/notifications'
                    aria-label='Notifications'
                    icon={<FaBell />}
                    variant='ghost'
                  />
                  <Badge
                    position='absolute'
                    top='-2px'
                    right='-2px'
                    colorScheme='red'
                    borderRadius='full'
                    size='xs'
                  >
                    3
                  </Badge>
                </Box>

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

                <Menu>
                  <MenuButton
                    as={Button}
                    rightIcon={<ChevronDownIcon />}
                    variant='ghost'
                    display={{ base: 'none', md: 'flex' }}
                  >
                    <Avatar
                      size='xs'
                      name={user?.name || 'User'}
                      src={user?.avatar}
                      mr={2}
                    />
                    <Text display={{ base: 'none', lg: 'block' }}>
                      {user?.name || 'Tài khoản'}
                    </Text>
                  </MenuButton>
                  <MenuList>
                    <MenuItem as={RouterLink} to='/account/profile'>
                      Tài khoản của tôi
                    </MenuItem>
                    <MenuItem as={RouterLink} to='/account/orders'>
                      Đơn hàng
                    </MenuItem>
                    <MenuDivider />
                    <MenuItem onClick={handleLogout}>Đăng xuất</MenuItem>
                  </MenuList>
                </Menu>
              </>
            ) : (
              <>
                <Button
                  as={RouterLink}
                  to='/login'
                  variant='ghost'
                  colorScheme='brand'
                  display={{ base: 'none', md: 'flex' }}
                >
                  Đăng nhập
                </Button>
                <Button as={RouterLink} to='/register' colorScheme='brand'>
                  Đăng ký
                </Button>
              </>
            )}
          </HStack>
        </Flex>
      </Container>

      {/* Mobile Drawer Menu */}
      <Drawer isOpen={isOpen} placement='left' onClose={onClose}>
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader>
            <Logo />
          </DrawerHeader>

          <DrawerBody>
            <VStack align='stretch' spacing={4}>
              {isAuthenticated ? (
                <Flex alignItems='center' p={2} mb={4}>
                  <Avatar
                    size='md'
                    name={user?.name}
                    src={user?.avatar}
                    mr={3}
                  />
                  <Box>
                    <Text fontWeight='bold'>{user?.name}</Text>
                    <Text fontSize='sm' color='gray.600'>
                      {user?.email}
                    </Text>
                  </Box>
                </Flex>
              ) : (
                <HStack justify='space-between' mb={4}>
                  <Button
                    as={RouterLink}
                    to='/login'
                    colorScheme='brand'
                    variant='outline'
                    flex='1'
                    onClick={onClose}
                  >
                    Đăng nhập
                  </Button>
                  <Button
                    as={RouterLink}
                    to='/register'
                    colorScheme='brand'
                    flex='1'
                    onClick={onClose}
                  >
                    Đăng ký
                  </Button>
                </HStack>
              )}

              <Box as='form' onSubmit={handleSearch} mb={4}>
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

              <Box as={RouterLink} to='/' py={2} onClick={onClose}>
                Trang chủ
              </Box>
              <Box as={RouterLink} to='/products' py={2} onClick={onClose}>
                Sản phẩm
              </Box>
              <Box as={RouterLink} to='/categories' py={2} onClick={onClose}>
                Danh mục
              </Box>
              <Box as={RouterLink} to='/promotions' py={2} onClick={onClose}>
                Khuyến mãi
              </Box>

              {isAuthenticated && (
                <>
                  <Box as='hr' my={2} borderColor='gray.200' />

                  <Box
                    as={RouterLink}
                    to='/account/profile'
                    py={2}
                    onClick={onClose}
                  >
                    Tài khoản của tôi
                  </Box>
                  <Box
                    as={RouterLink}
                    to='/account/orders'
                    py={2}
                    onClick={onClose}
                  >
                    Đơn hàng
                  </Box>
                  <Box as={RouterLink} to='/favorites' py={2} onClick={onClose}>
                    Sản phẩm yêu thích
                  </Box>
                  <Box as={RouterLink} to='/cart' py={2} onClick={onClose}>
                    Giỏ hàng
                  </Box>
                  <Box
                    as={RouterLink}
                    to='/notifications'
                    py={2}
                    onClick={onClose}
                  >
                    Thông báo
                  </Box>

                  <Box as='hr' my={2} borderColor='gray.200' />

                  <Box
                    as='button'
                    py={2}
                    color='red.500'
                    textAlign='left'
                    onClick={() => {
                      handleLogout();
                      onClose();
                    }}
                  >
                    Đăng xuất
                  </Box>
                </>
              )}
            </VStack>
          </DrawerBody>
        </DrawerContent>
      </Drawer>
    </Box>
  );
};

export default Header;
