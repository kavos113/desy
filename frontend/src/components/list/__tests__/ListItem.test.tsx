import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import ListItem from "../ListItem";
import { domain } from "../../../../wailsjs/go/models";

describe("ListItem", () => {
  const createSummary = () =>
    domain.LectureSummary.createFrom({
      ID: 1,
      University: "東京工業大学",
      Title: "アルゴリズム",
      Department: "情報理工学院",
      Code: "ABC123",
      Level: 3,
      Year: 2024,
      Timetables: [
        {
          DayOfWeek: "monday",
          Period: 1,
          Room: { Name: "W1" },
          Semester: "spring",
        },
        {
          DayOfWeek: "monday",
          Period: 2,
          Room: { Name: "W1" },
          Semester: "spring",
        },
      ],
      Teachers: [
        {
          ID: 1,
          Name: "田中太郎",
          Url: "https://example.com",
        },
      ],
    });

  it("講義情報を表示し、クリック時にIDを通知する", () => {
    const handleClick = vi.fn();
    render(<ListItem item={createSummary()} onClick={handleClick} />);

    expect(screen.getByText("東京工業大学")).toBeInTheDocument();
    expect(screen.getByText("ABC123")).toBeInTheDocument();
    expect(screen.getByText("アルゴリズム")).toBeInTheDocument();
    expect(screen.getByText("田中太郎")).toBeInTheDocument();
  expect(screen.getByText("月1-2(W1)")).toBeInTheDocument();

    fireEvent.click(screen.getByText("アルゴリズム"));
    expect(handleClick).toHaveBeenCalledWith(1);
  });
});
