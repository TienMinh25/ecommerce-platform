import React from 'react';
import { Link } from 'react-router-dom';
import { Box, Flex, Text } from '@chakra-ui/react';

// Logo sẽ có 3 chế độ: full (icon + text), icon (chỉ biểu tượng), text (chỉ text)
// Và 3 style: minimal, modern, shopping
const Logo = ({ mode = 'full', size = 'md', style = 'minimal', ...props }) => {
  const sizeValues = {
    sm: { icon: 30, fontSize: '18px' },
    md: { icon: 40, fontSize: '24px' },
    lg: { icon: 50, fontSize: '30px' },
  };

  const { icon: iconSize, fontSize } = sizeValues[size] || sizeValues.md;

  const colors = {
    minimal: { primary: '#3B82F6', text: '#1E293B' },
    modern: { primary: '#4F46E5', text: '#111827' },
    shopping: { primary: '#FF385C', text: '#333333' },
  };

  const { primary, text } = colors[style] || colors.minimal;

  return (
    <Link to='/'>
      <Flex align='center' {...props}>
        {(mode === 'full' || mode === 'icon') && (
          <LogoIcon
            boxSize={iconSize}
            mr={mode === 'full' ? 3 : 0}
            style={style}
            primaryColor={primary}
          />
        )}

        {(mode === 'full' || mode === 'text') && (
          <Text
            fontWeight='800'
            fontSize={fontSize}
            letterSpacing='wide'
            color={text}
            textTransform={style === 'minimal' ? 'capitalize' : 'uppercase'}
          >
            {style === 'minimal' ? 'Minh Plaza' : 'MINH PLAZA'}
          </Text>
        )}
      </Flex>
    </Link>
  );
};

const LogoIcon = ({ style = 'minimal', primaryColor, ...props }) => {
  const renderIcon = () => {
    switch (style) {
      case 'modern':
        return (
          <>
            {/* Circle background */}
            <circle cx='30' cy='30' r='30' fill={primaryColor} />

            {/* MP letters in negative space */}
            <path
              d='M15,15 L15,45 L25,45 L25,35 L35,35 L35,45 L45,45 L45,15 L35,15 L35,25 L25,25 L25,15 Z'
              fill='white'
            />
          </>
        );

      case 'shopping':
        return (
          <>
            {/* Shopping bag background shape */}
            <path
              d='M0,15 C0,10 5,0 25,0 C45,0 50,10 50,15 L60,60 L-10,60 L0,15 Z'
              fill={primaryColor}
            />

            {/* M shape in white */}
            <path
              d='M10,15 L10,45 L25,25 L40,45 L40,15'
              stroke='white'
              strokeWidth='5'
              fill='none'
              strokeLinecap='round'
              strokeLinejoin='round'
            />
          </>
        );

      case 'minimal':
      default:
        return (
          <>
            {/* Simple square */}
            <rect
              x='0'
              y='0'
              width='50'
              height='50'
              rx='10'
              fill={primaryColor}
            />

            {/* M letter */}
            <path
              d='M10,10 L10,40 L25,25 L40,40 L40,10'
              stroke='white'
              strokeWidth='6'
              fill='none'
              strokeLinecap='round'
              strokeLinejoin='round'
            />
          </>
        );
    }
  };

  return (
    <Box {...props}>
      <svg
        xmlns='http://www.w3.org/2000/svg'
        viewBox='0 0 60 60'
        width='100%'
        height='100%'
      >
        {renderIcon()}
      </svg>
    </Box>
  );
};

export default Logo;
