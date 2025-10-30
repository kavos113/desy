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

const DAY_ORDER = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
];

type TimeTableGroup = {
  dayKey: string;
  dayLabel: string;
  roomLabel: string;
  order: number;
  periods: Set<number>;
};

type PeriodRange = {
  start: number;
  end: number;
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

  const groups = new Map<string, TimeTableGroup>();
  const fallbacks: string[] = [];

  timetables.forEach((timetable) => {
    if (!timetable) {
      return;
    }

    const rawDay = timetable.DayOfWeek ?? "";
    const dayKey = rawDay.toString().toLowerCase();
    const dayLabel = DAY_OF_WEEK_LABELS[dayKey] ?? rawDay;
    const roomLabel = timetable.Room?.Name?.trim() ?? "";
    const periodValue = Number(timetable.Period);
    const hasPeriod = Number.isFinite(periodValue) && periodValue > 0;

    const fallback = buildFallbackLabel(
      dayLabel,
      timetable.Period,
      timetable.Room?.Name
    );

    if (!dayLabel || !hasPeriod) {
      if (fallback) {
        fallbacks.push(fallback);
      }
      return;
    }

    const key = `${dayKey}::${roomLabel}`;
    if (!groups.has(key)) {
      const orderIndex = DAY_ORDER.indexOf(dayKey);
      groups.set(key, {
        dayKey,
        dayLabel,
        roomLabel,
        order: orderIndex === -1 ? Number.MAX_SAFE_INTEGER : orderIndex,
        periods: new Set<number>(),
      });
    }

    groups.get(key)!.periods.add(periodValue);
  });

  const formatted: string[] = [];

  Array.from(groups.values())
    .sort((left, right) => {
      if (left.order !== right.order) {
        return left.order - right.order;
      }

      if (left.roomLabel === right.roomLabel) {
        return left.dayLabel.localeCompare(right.dayLabel, "ja");
      }

      if (!left.roomLabel) {
        return -1;
      }
      if (!right.roomLabel) {
        return 1;
      }
      return left.roomLabel.localeCompare(right.roomLabel, "ja");
    })
    .forEach((group) => {
      const periods = Array.from(group.periods).sort((a, b) => a - b);
      const ranges = compressPeriods(periods);
      const roomSuffix = group.roomLabel ? `(${group.roomLabel})` : "";

      ranges.forEach((range) => {
        const periodLabel =
          range.start === range.end
            ? `${range.start}`
            : `${range.start}-${range.end}`;
        formatted.push(`${group.dayLabel}${periodLabel}${roomSuffix}`);
      });
    });

  const uniqueFallbacks = fallbacks.filter(Boolean);

  const result = [...formatted, ...uniqueFallbacks];
  return Array.from(new Set(result)).join(", ");
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

export function splitIntoLines(value: string | undefined | null): string[] {
  if (!value) {
    return [];
  }

  return value
    .split(/<br\s*\/?>|\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
}

function compressPeriods(periods: number[]): PeriodRange[] {
  if (periods.length === 0) {
    return [];
  }

  const ranges: PeriodRange[] = [];
  let start = periods[0];
  let end = start;

  for (let index = 1; index < periods.length; index += 1) {
    const current = periods[index];
    if (current === end + 1) {
      end = current;
      continue;
    }

    ranges.push({ start, end });
    start = current;
    end = current;
  }

  ranges.push({ start, end });
  return ranges;
}

function buildFallbackLabel(
  dayLabel: string,
  period: number | undefined,
  roomName: string | undefined
): string {
  const dayPart = dayLabel ?? "";
  const periodPart =
    Number.isFinite(Number(period)) && Number(period) > 0
      ? `${Number(period)}`
      : "";
  const roomPart = roomName ? `(${roomName})` : "";

  if (!dayPart && !periodPart && !roomPart) {
    return "";
  }

  return `${dayPart}${periodPart}${roomPart}`;
}
