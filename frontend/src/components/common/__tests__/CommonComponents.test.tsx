import { fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import CheckBox from "../CheckBox";
import ComboBox from "../ComboBox";
import SearchBox from "../SearchBox";
import SimpleButton from "../SimpleButton";

describe("CheckBox", () => {
  it("toggles state and emits parsed index", async () => {
    const user = userEvent.setup();
    const handleCheckItem = vi.fn();

    render(<CheckBox checkboxId="grade-1" content="学士1年" onCheckItem={handleCheckItem} />);

    const checkbox = screen.getByRole("checkbox");
    await user.click(checkbox);
    expect(handleCheckItem).toHaveBeenLastCalledWith(1, true);

    await user.click(checkbox);
    expect(handleCheckItem).toHaveBeenLastCalledWith(1, false);
  });
});

describe("SearchBox", () => {
  it("notifies text changes", async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();

    render(<SearchBox placeholder="講義名" onChange={handleChange} />);

    const input = screen.getByPlaceholderText("講義名");
    await user.type(input, "AB");

    expect(handleChange).toHaveBeenNthCalledWith(1, "A");
    expect(handleChange).toHaveBeenNthCalledWith(2, "AB");
  });
});

describe("SimpleButton", () => {
  it("fires click handler", async () => {
    const user = userEvent.setup();
    const handleClick = vi.fn();

    render(<SimpleButton text="検索" onClick={handleClick} />);

    const button = screen.getByRole("button", { name: "検索" });
    await user.click(button);

    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});

describe("ComboBox", () => {
  it("toggles selection for simple lists", async () => {
    const user = userEvent.setup();
    const handleSelect = vi.fn();

    render(<ComboBox items={["A", "B"]} onSelectItem={handleSelect} />);

    const itemA = screen.getByText("A");
    await user.click(itemA);
    expect(handleSelect).toHaveBeenLastCalledWith(["A"]);

    await user.click(itemA);
    expect(handleSelect).toHaveBeenLastCalledWith([]);
  });

  it("supports nested menu selections", async () => {
    const user = userEvent.setup();
    const handleSelect = vi.fn();

    render(
      <ComboBox
        items={{
          Parent: ["Child"],
        }}
        onSelectItem={handleSelect}
      />
    );

    const parent = screen.getByText("Parent");
    const listItem = parent.closest("li");
    if (!listItem) {
      throw new Error("Menu item wrapper not found");
    }

    fireEvent.mouseEnter(listItem);

    const child = await screen.findByText("Child");
    expect(child).toBeVisible();

    await user.click(child);
    expect(handleSelect).toHaveBeenLastCalledWith(["Child"]);
  });
});
