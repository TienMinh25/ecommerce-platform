import { createContext, useState, useContext, useEffect } from 'react';
import cartService from '../services/cartService';
import { useToast } from '@chakra-ui/react';

const CartContext = createContext();

export const CartProvider = ({ children }) => {
    const [cartItems, setCartItems] = useState([]);
    const [selectedItems, setSelectedItems] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const toast = useToast();

    const fetchCartItems = async () => {
        try {
            setIsLoading(true);
            const response = await cartService.getCartItems();
            const items = Array.isArray(response.data) ? response.data : [];
            setCartItems(items);
        } catch (error) {
            console.error('Error fetching cart items:', error);
        } finally {
            setIsLoading(false);
        }
    };

    // Refresh cart items
    const refreshCart = () => {
        fetchCartItems();
    };

    useEffect(() => {
        fetchCartItems();
    }, []);

    // Add item to cart
    const addToCart = async (cartItem) => {
        try {
            setIsLoading(true);
            await cartService.addToCart(cartItem);
            await fetchCartItems(); // Refresh cart after adding
            return true;
        } catch (error) {
            console.error('Error adding item to cart:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể thêm sản phẩm vào giỏ hàng',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            return false;
        } finally {
            setIsLoading(false);
        }
    };

    // Update cart item
    const updateCartItem = async (cartItemId, data) => {
        try {
            await cartService.updateCartItem(cartItemId, data);

            // Update local state
            if (data.quantity > 0) {
                setCartItems(prevItems =>
                    prevItems.map(item =>
                        item.cart_item_id === cartItemId ? { ...item, quantity: data.quantity } : item
                    )
                );
            } else {
                // If quantity is 0, remove from cart
                setCartItems(prevItems => prevItems.filter(item => item.cart_item_id !== cartItemId));
                setSelectedItems(prev => prev.filter(id => id !== cartItemId));
            }

            return true;
        } catch (error) {
            console.error('Error updating cart item:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể cập nhật số lượng sản phẩm',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            return false;
        }
    };

    // Delete cart items
    const deleteCartItems = async (cartItemIds) => {
        try {
            await cartService.deleteCartItems(cartItemIds);

            // Update local state
            setCartItems(prevItems =>
                prevItems.filter(item => !cartItemIds.includes(item.cart_item_id))
            );
            setSelectedItems(prev =>
                prev.filter(id => !cartItemIds.includes(id))
            );

            toast({
                title: 'Thành công',
                description: 'Đã xóa sản phẩm khỏi giỏ hàng',
                status: 'success',
                duration: 3000,
                isClosable: true,
            });

            return true;
        } catch (error) {
            console.error('Error deleting cart items:', error);
            toast({
                title: 'Lỗi',
                description: 'Không thể xóa sản phẩm khỏi giỏ hàng',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            return false;
        }
    };

    const toggleSelectItem = (cartItemId) => {
        if (selectedItems.includes(cartItemId)) {
            setSelectedItems(prev => prev.filter(id => id !== cartItemId));
        } else {
            setSelectedItems(prev => [...prev, cartItemId]);
        }
    };

    const toggleSelectAll = () => {
        if (selectedItems.length === cartItems.length) {
            setSelectedItems([]);
        } else {
            setSelectedItems(cartItems.map(item => item.cart_item_id));
        }
    };

    // Format price in VND
    const formatPrice = (price) => {
        return new Intl.NumberFormat('vi-VN', {
            style: 'currency',
            currency: 'VND',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(price);
    };

    // Calculate total price for selected items
    const calculateTotal = () => {
        return cartItems
            .filter(item => selectedItems.includes(item.cart_item_id))
            .reduce((total, item) => {
                const price = item.discount_price > 0 ? item.discount_price : item.price;
                return total + (price * item.quantity);
            }, 0);
    };

    return (
        <CartContext.Provider
            value={{
                cartItems,
                setCartItems,
                selectedItems,
                setSelectedItems,
                isLoading,
                refreshCart,
                addToCart,
                updateCartItem,
                deleteCartItems,
                toggleSelectItem,
                toggleSelectAll,
                formatPrice,
                calculateTotal,
                cartCount: cartItems.length
            }}
        >
            {children}
        </CartContext.Provider>
    );
};

export const useCart = () => useContext(CartContext);