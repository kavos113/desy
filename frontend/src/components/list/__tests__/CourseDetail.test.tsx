import { render, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import CourseDetail from "../CourseDetail";
import { domain } from "../../../../wailsjs/go/models";

describe("CourseDetail", () => {
  const createLecture = (id: number) =>
    new domain.Lecture({
      ID: id,
      Title: `Lecture ${id}`,
      EnglishTitle: "",
      Department: "",
      LectureType: "",
      Timetables: [],
      Teachers: [],
      LecturePlans: [],
      Keywords: [],
      RelatedCourseCodes: [],
      RelatedCourses: [],
    });

  it("スクロール位置を先頭に戻す", async () => {
    const detailPanel = document.createElement("div");
    detailPanel.className = "detail-panel";
  const scrollToMock = vi.fn();
  (detailPanel as any).scrollTo = scrollToMock;
    document.body.appendChild(detailPanel);

    const lecture1 = createLecture(1);
    const lecture2 = createLecture(2);

    const renderResult = render(<CourseDetail lecture={lecture1} relatedCourses={[]} />, {
      container: detailPanel,
    });

    expect(scrollToMock).toHaveBeenCalledWith({ top: 0, behavior: "auto" });
    scrollToMock.mockClear();

    const wrapper = detailPanel.querySelector<HTMLDivElement>(".course-detail-wrapper");
    expect(wrapper).not.toBeNull();
    if (!wrapper) {
      throw new Error("wrapper not found");
    }

    wrapper.scrollTop = 200;

    renderResult.rerender(<CourseDetail lecture={lecture2} relatedCourses={[]} />);

    await waitFor(() => {
      expect(scrollToMock).toHaveBeenCalledWith({ top: 0, behavior: "auto" });
    });
    expect(wrapper.scrollTop).toBe(0);

    renderResult.unmount();
    document.body.removeChild(detailPanel);
  });
});
