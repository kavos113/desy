import { useMemo } from "react";
import { domain } from "../../../wailsjs/go/models";
import {
  formatRelatedCourses,
  formatSemesters,
  formatTeachers,
  formatTimetables,
  splitIntoLines,
} from "./utils";
import "./list.css";

type CourseDetailProps = {
  lecture?: domain.Lecture | null;
};

const CourseDetail = ({ lecture }: CourseDetailProps) => {
  const teachers = useMemo(() => lecture?.Teachers ?? [], [lecture?.Teachers]);
  const teacherAnchors = teachers.map((teacher) => (
    <a key={teacher.ID} href={teacher.Url ?? undefined} className="lecturerText">
      {teacher.Name}
    </a>
  ));

  const timetableText = useMemo(() => formatTimetables(lecture?.Timetables), [lecture?.Timetables]);
  const semesterText = useMemo(() => formatSemesters(lecture?.Timetables), [lecture?.Timetables]);
  const relatedCourses = useMemo(() => formatRelatedCourses(lecture?.RelatedCourses), [lecture?.RelatedCourses]);

  const keywords = lecture?.Keywords ?? [];

  const plans = lecture?.LecturePlans ?? [];

  const renderParagraphs = (value?: string | null) => {
    return splitIntoLines(value).map((line, index) => (
      <p key={`${line}-${index}`} className="course-detail-text">
        {line}
      </p>
    ));
  };

  const renderList = (value?: string | null) => {
    const items = splitIntoLines(value);
    return items.map((item, index) => <li key={`${item}-${index}`}>{item}</li>);
  };

  if (!lecture) {
    return (
      <div className="course-detail-wrapper">
        <p>講義の詳細が選択されていません。</p>
      </div>
    );
  }

  return (
    <div className="course-detail-wrapper">
      <div className="course-detail-title">
        <p className="course-title-ja">{lecture.Title}</p>
        <p className="course-title-en">{lecture.EnglishTitle}</p>
      </div>

      <div className="course-detail-grid">
        <dl className="course-detail-row">
          <dt>開講元</dt>
          <dd>{lecture.Department}</dd>
        </dl>
        <dl className="course-detail-row">
          <dt>担当教員</dt>
          <dd>{teacherAnchors}</dd>
        </dl>
        <dl className="course-detail-row">
          <dt>授業形態</dt>
          <dd>{lecture.LectureType}</dd>
        </dl>
        <dl className="course-detail-row">
          <dt>曜日・時限(講義室)</dt>
          <dd>{timetableText}</dd>
        </dl>
        <dl className="course-detail-row">
          <dt>開講時期</dt>
          <dd>
            {lecture.Year ? `${lecture.Year}年` : null}
            {semesterText ? ` ${semesterText}` : null}
          </dd>
        </dl>
      </div>

      <div className="course-detail-row">
        <dl className="course-detail-row-half">
          <dt>科目コード</dt>
          <dd>{lecture.Code}</dd>
        </dl>
        <dl className="course-detail-row-half">
          <dt>単位数</dt>
          <dd>{lecture.Credit}</dd>
        </dl>
      </div>

      <div className="course-detail-row">
        <dl className="course-detail-row-half">
          <dt>言語</dt>
          <dd>{lecture.Language}</dd>
        </dl>
        <dl className="course-detail-row-half">
          <dt>シラバスURL</dt>
          <dd>
            {lecture.Url ? (
              <a href={lecture.Url} target="_blank" rel="noreferrer">
                {lecture.Url}
              </a>
            ) : null}
          </dd>
        </dl>
      </div>

      <div className="course-detail-section">
        <h3>講義の概要とねらい</h3>
        {renderParagraphs(lecture.Abstract)}
      </div>

      <div className="course-detail-section">
        <h3>授業計画・課題</h3>
        <table className="course-plan-table">
          <thead>
            <tr>
              <th>回</th>
              <th>授業計画</th>
              <th>課題</th>
            </tr>
          </thead>
          <tbody>
            {plans.map((plan) => (
              <tr key={plan.Count}>
                <td className="course-plan-count">第{plan.Count}回</td>
                <td>{plan.Plan}</td>
                <td>{plan.Assignment}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="course-detail-section">
        <h3>到達目標</h3>
        {renderParagraphs(lecture.Goal)}
      </div>

      {lecture.Experience ? (
        <div className="course-detail-section">
          <h3>実務経験のある教員による授業</h3>
          <p className="course-detail-text">あり</p>
        </div>
      ) : null}

      <div className="course-detail-section">
        <h3>キーワード</h3>
        <p className="course-detail-text">{keywords.join(", ")}</p>
      </div>

      <div className="course-detail-section">
        <h3>教科書</h3>
        <ul className="course-detail-reference">{renderList(lecture.Textbook)}</ul>
      </div>

      <div className="course-detail-section">
        <h3>参考書・講義資料等</h3>
        <ul className="course-detail-reference">{renderList(lecture.ReferenceBook)}</ul>
      </div>

      <div className="course-detail-section">
        <h3>授業の進め方</h3>
        {renderParagraphs(lecture.Flow)}
      </div>

      <div className="course-detail-section">
        <h3>授業時間外学修（予習・復習等）</h3>
        {renderParagraphs(lecture.OutOfClassWork)}
      </div>

      <div className="course-detail-section">
        <h3>成績評価の基準及び方法</h3>
        {renderParagraphs(lecture.Assessment)}
      </div>

      <div className="course-detail-section">
        <h3>関連する科目</h3>
        <p className="course-detail-text">{relatedCourses}</p>
      </div>

      <div className="course-detail-section">
        <h3>履修の条件</h3>
        {renderParagraphs(lecture.Prerequisite)}
      </div>

      <div className="course-detail-section">
        <h3>その他</h3>
        {renderParagraphs(lecture.Note)}
      </div>

      <div className="course-detail-section">
        <h3>連絡先</h3>
        {renderParagraphs(lecture.Contact)}
      </div>

      <div className="course-detail-section">
        <h3>オフィスアワー</h3>
        {renderParagraphs(lecture.OfficeHours)}
      </div>
    </div>
  );
};

export default CourseDetail;
