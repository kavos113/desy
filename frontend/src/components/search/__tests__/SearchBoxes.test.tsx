import { act, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import SearchBoxes from "../SearchBoxes";

describe("SearchBoxes", () => {
  it("選択したコンボボックスの値を通知する", async () => {
    const handleSelect = vi.fn();
    const user = userEvent.setup();
    const { container } = render(<SearchBoxes onSelectMenuItem={handleSelect} />);

    const selects = container.querySelectorAll("select");
    const universitySelect = selects[0] as HTMLSelectElement;
    const departmentSelect = selects[1] as HTMLSelectElement;

    await act(async () => {
      await user.selectOptions(universitySelect, ["東京工業大学"]);
      await user.selectOptions(departmentSelect, ["理学院", "工学院"]);
    });

    await waitFor(() => {
      const universityCall = handleSelect.mock.calls.find(
        ([key, items]) => key === "university" && (items as string[]).length > 0
      );
      const departmentCall = handleSelect.mock.calls
        .filter(([key]) => key === "department")
        .at(-1);

      expect(universityCall?.[1]).toEqual(["東京工業大学"]);
      expect(departmentCall?.[1]).toEqual(expect.arrayContaining(["理学院", "工学院"]));
      expect((departmentCall?.[1] as string[] | undefined)?.length).toBe(2);
    });
  });

  it("検索ボックスの入力値を通知する", async () => {
    const handleChange = vi.fn();
    const user = userEvent.setup();
    render(<SearchBoxes onChangeSearchBox={handleChange} />);

    const inputs = screen.getAllByRole("textbox");
    const titleInput = inputs[0];
    const lecturerInput = inputs[1];

    await act(async () => {
      await user.type(titleInput, "データサイエンス");
      await user.type(lecturerInput, "田中");
    });

    await waitFor(() => {
      const lastCall = handleChange.mock.calls.at(-1);
      expect(lastCall).toEqual(["データサイエンス", "田中"]);
    });
  });
});
