import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import Timetable from "../Timetable";

describe("Timetable", () => {
  it("セルをクリックすると選択状態がトグルする", async () => {
    const handleSelect = vi.fn();
    render(<Timetable onCheckItem={handleSelect} />);

    expect(handleSelect).toHaveBeenCalledWith([]);
    handleSelect.mockClear();

    const cells = screen.getAllByRole("cell");
    fireEvent.click(cells[0]);

    await waitFor(() => {
      expect(handleSelect).toHaveBeenLastCalledWith([{ day: "月", period: "1" }]);
    });

    fireEvent.click(cells[0]);

    await waitFor(() => {
      expect(handleSelect).toHaveBeenLastCalledWith([]);
    });
  });

  it("曜日ヘッダーをクリックすると列全体がトグルする", async () => {
    const handleSelect = vi.fn();
    render(<Timetable onCheckItem={handleSelect} />);

    handleSelect.mockClear();

    const mondayHeader = screen.getByText("月");
    fireEvent.click(mondayHeader);

    await waitFor(() => {
      expect(handleSelect).toHaveBeenLastCalledWith([
        { day: "月", period: "1" },
        { day: "月", period: "2" },
        { day: "月", period: "3" },
        { day: "月", period: "4" },
        { day: "月", period: "5" },
      ]);
    });
  });
});
