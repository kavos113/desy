import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import CheckBoxes from "../CheckBoxes";

describe("CheckBoxes", () => {
  it("選択状態をコールバックへ伝播する", async () => {
    const handleCheck = vi.fn();
    render(
      <CheckBoxes
        checkboxId="grade"
        contents={["A", "B", "C"]}
        onCheckItem={handleCheck}
      />
    );

    expect(handleCheck).toHaveBeenCalledWith([]);
    handleCheck.mockClear();

    const checkboxA = screen.getByLabelText("A");
    fireEvent.click(checkboxA);

    await waitFor(() => {
      expect(handleCheck).toHaveBeenLastCalledWith(["A"]);
    });

    const checkboxB = screen.getByLabelText("B");
    fireEvent.click(checkboxB);

    await waitFor(() => {
      expect(handleCheck).toHaveBeenLastCalledWith(["A", "B"]);
    });

    fireEvent.click(checkboxA);

    await waitFor(() => {
      expect(handleCheck).toHaveBeenLastCalledWith(["B"]);
    });
  });
});
