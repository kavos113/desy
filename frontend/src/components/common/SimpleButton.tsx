import { ButtonHTMLAttributes } from 'react';
import './common.css';

type SimpleButtonProps = {
  text: string;
  className?: string;
  type?: 'button' | 'submit' | 'reset';
} & Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'children' | 'className' | 'type'>;

const SimpleButton = ({ text, className, type = 'button', ...buttonProps }: SimpleButtonProps) => {
  const mergedClassName = ['simple-button', className].filter(Boolean).join(' ');

  return (
    <button type={type} className={mergedClassName} {...buttonProps}>
      <span className="simple-button-text">{text}</span>
    </button>
  );
};

export default SimpleButton;
