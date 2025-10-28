import { ButtonHTMLAttributes } from "react";
import "./common.css";

type SimpleButtonProps = {
  text: string;
} & Omit<ButtonHTMLAttributes<HTMLButtonElement>, "type">;

const SimpleButton = ({ text, className, ...rest }: SimpleButtonProps) => {
  const buttonClass = ["simple-button", className].filter(Boolean).join(" ");

  return (
    <button type="button" className={buttonClass} {...rest}>
      <span className="simple-button-text">{text}</span>
    </button>
  );
};

export default SimpleButton;
