import SimpleButton from "./SimpleButton";

export default {
  title: "common/SimpleButton",
  component: SimpleButton,
  argTypes: {
    text: {control: "text"},
  },
};

export const Default = {
  args: {
    primary: true,
    text: "Click Me",
  },
};