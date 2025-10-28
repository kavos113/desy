import { domain } from "../../../wailsjs/go/models";

const DAY_OF_WEEK_LABELS: Record<string, string> = {
  monday: "月",
  tuesday: "火",
  wednesday: "水",
  thursday: "木",
  friday: "金",
  saturday: "土",
  sunday: "日",
};

const SEMESTER_LABELS: Record<string, string> = {
  spring: "春学期",
  summer: "夏学期",
  fall: "秋学期",
  autumn: "秋学期",
  winter: "冬学期",
};

export function formatTeachers(teachers: domain.Teacher[] | undefined): string {
  if (!teachers || teachers.length === 0) {
    return "";
  }
  return teachers.map((teacher) => teacher.Name).join(", ");
}

export function formatTimetables(
  timetables: domain.TimeTable[] | undefined
): string {
  if (!timetables || timetables.length === 0) {
    return "";
  }

  const formatted = timetables.map((timetable) => {
    const day =
      DAY_OF_WEEK_LABELS[timetable.DayOfWeek?.toLowerCase() ?? ""] ??
      timetable.DayOfWeek;
    const period = timetable.Period ? `${timetable.Period}` : "";
    const room = timetable.Room?.Name ? `(${timetable.Room.Name})` : "";
    return `${day}${period}${room}`;
  });

  return Array.from(new Set(formatted)).join(", ");
}

export function formatSemesters(
  timetables: domain.TimeTable[] | undefined
): string {
  if (!timetables || timetables.length === 0) {
    return "";
  }

  const labels = timetables
    .map(
      (timetable) =>
        SEMESTER_LABELS[timetable.Semester?.toLowerCase() ?? ""] ??
        timetable.Semester
    )
    .filter((label): label is string => Boolean(label));

  if (labels.length === 0) {
    return "";
  }

  return Array.from(new Set(labels)).join(", ");
}

export function formatRelatedCourses(courses: number[] | undefined): string {
  if (!courses || courses.length === 0) {
    return "";
  }

  return courses.map((course) => course.toString()).join(", ");
}

export function splitIntoLines(value: string | undefined | null): string[] {
  if (!value) {
    return [];
  }

  return value
    .split(/<br\s*\/?>|\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
}
