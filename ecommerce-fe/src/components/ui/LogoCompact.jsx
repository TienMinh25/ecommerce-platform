import React from 'react';
import { Link } from 'react-router-dom';
import { Box } from '@chakra-ui/react';

// CompactLogo component with three sizes and three styles
// This version doesn't include the text part, only the icon
const LogoCompact = ({ size = 'md', style = 'minimal', ...props }) => {
  const sizeValues = {
    sm: 30,
    md: 40,
    lg: 50,
  };

  const iconSize = sizeValues[size] || sizeValues.md;

  const colors = {
    minimal: { primary: '#3B82F6' },
    modern: { primary: '#4F46E5' },
    shopping: { primary: '#FF385C' },
  };

  const { primary } = colors[style] || colors.minimal;

  return (
    <Link to='/'>
      <Box boxSize={iconSize} {...props}>
        <LogoIcon 
          style={style} 
          primaryColor={primary} 
        />
      </Box>
    </Link>
  );
};

const LogoIcon = ({ style = 'minimal', primaryColor }) => {
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
    <svg
      xmlns='http://www.w3.org/2000/svg'
      viewBox='0 0 60 60'
      width='100%'
      height='100%'
    >
      {renderIcon()}
    </svg>
  );
};

export default LogoCompact;