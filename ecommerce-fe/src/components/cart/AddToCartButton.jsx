import React, { useState } from 'react';
import { Button, Icon, useToast } from '@chakra-ui/react';
import { FaShoppingCart } from 'react-icons/fa';
import {useCart} from "../../context/CartContext.jsx";

const AddToCartButton = ({ productId, variantId, quantity = 1, onSuccess, ...props }) => {
    const [isLoading, setIsLoading] = useState(false);
    const { addToCart } = useCart();
    const toast = useToast();

    const handleAddToCart = async () => {
        if (!productId || !variantId) {
            toast({
                title: 'Lỗi',
                description: 'Thiếu thông tin sản phẩm cần thiết',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
            return;
        }

        try {
            setIsLoading(true);
            const success = await addToCart({
                product_id: productId,
                product_variant_id: variantId,
                quantity: quantity
            });

            if (success && onSuccess && typeof onSuccess === 'function') {
                onSuccess();
            }
        } catch (error) {
            toast({
                title: 'Lỗi',
                description: 'Không thể thêm sản phẩm vào giỏ hàng',
                status: 'error',
                duration: 3000,
                isClosable: true,
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Button
            leftIcon={<Icon as={FaShoppingCart} />}
            onClick={handleAddToCart}
            isLoading={isLoading}
            {...props}
        >
            {props.children || 'Thêm vào giỏ'}
        </Button>
    );
};

export default AddToCartButton;