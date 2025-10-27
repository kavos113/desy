import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import SearchConditions from "../SearchConditions";

describe("SearchConditions", () => {
  it("チェックボックスの選択を通知する", async () => {
    const handleCheck = vi.fn();
    render(<SearchConditions onCheckItem={handleCheck} />);

    const gradeCheckbox = screen.getByLabelText("学士1年");
    const quarterCheckbox = screen.getByLabelText("1Q");

    await userEvent.click(gradeCheckbox);
    await userEvent.click(quarterCheckbox);

    await waitFor(() => {
      expect(handleCheck).toHaveBeenCalledWith("grade", ["学士1年"]);
      expect(handleCheck).toHaveBeenCalledWith("quarter", ["1Q"]);
    });
  });

  it("時間割の選択を通知する", async () => {
    const handleTimetable = vi.fn();
    const { container } = render(<SearchConditions onTimetableChange={handleTimetable} />);

    const firstRowCells = container.querySelectorAll<HTMLTableCellElement>("tbody tr:first-child td");
    const firstCell = firstRowCells[0];
    if (!firstCell) {
      throw new Error("時間割セルが見つかりません");
    }

    await userEvent.click(firstCell);

    await waitFor(() => {
      expect(handleTimetable).toHaveBeenCalledWith([
        {
          day: "月",
          period: "1",
        },
      ]);
    });
  });
});
