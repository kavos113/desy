import "./button.css"

type SimpleButtonProps = {
  text: string;
};

const SimpleButton = ({ text }: SimpleButtonProps) => {
  return (
    <button type="button" className="simple-button">
      <span className="simple-button-text">{text}</span>
    </button>
  );
};

export default SimpleButton;