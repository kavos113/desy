import './list.css';
import { domain } from '../../../wailsjs/go/models';
import { LEVEL_LABELS, LECTURE_TYPE_LABELS } from '../../constants';

const DAY_OF_WEEK_LABELS: Record<string, string> = {
  monday: '月',
  tuesday: '火',
  wednesday: '水',
  thursday: '木',
  friday: '金',
  saturday: '土',
  sunday: '日'
};

const SEMESTER_LABELS: Record<string, string> = {
  spring: '春学期',
  summer: '夏学期',
  fall: '秋学期',
  winter: '冬学期'
};

const formatTimetable = (timetable: domain.TimeTable): string => {
  const semesterKey = timetable.Semester?.toLowerCase?.() ?? '';
  const semester = SEMESTER_LABELS[semesterKey] ?? timetable.Semester ?? '';
  const dayKey = timetable.DayOfWeek?.toLowerCase?.() ?? '';
  const day = DAY_OF_WEEK_LABELS[dayKey] ?? timetable.DayOfWeek ?? '';
  const period = timetable.Period ? `${timetable.Period}限` : '';
  const room = timetable.Room?.Name ? `(${timetable.Room.Name})` : '';
  const parts = [semester, day ? `${day}曜` : '', period].filter(Boolean);
  return parts.join(' ') + room;
};

const formatTeachers = (teachers: domain.Teacher[] | undefined): string =>
  (teachers ?? []).map((teacher) => teacher.Name).join(', ') || '-';

const formatLectureType = (lectureType: string | undefined): string => {
  if (!lectureType) {
    return '未設定';
  }
  const key = lectureType.toLowerCase();
  return LECTURE_TYPE_LABELS[key] ?? lectureType;
};

const formatLevelLabel = (level: number | undefined): string => {
  if (!level) {
    return '-';
  }
  return LEVEL_LABELS[level] ?? `${level}`;
};

const splitByBreak = (text: string | undefined): string[] => {
  if (!text) {
    return [];
  }
  return text
    .split(/<br\s*\/?>(?:\r?\n)?|\r?\n/g)
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
};

type CourseDetailPanelProps = {
  lecture: domain.Lecture | null;
  open: boolean;
  loading?: boolean;
  errorMessage?: string | null;
};

const CourseDetailPanel = ({ lecture, open, loading, errorMessage }: CourseDetailPanelProps) => {
  if (!open) {
    return null;
  }

  return (
    <aside
      className={`detail-panel${open ? ' detail-panel--open' : ''}`}
      role="dialog"
      aria-modal="true"
      aria-labelledby="lecture-detail-title"
    >
        <div className="detail-panel__inner">
          {loading ? (
            <p className="detail-loading">詳細を読み込み中です...</p>
          ) : errorMessage ? (
            <p className="detail-error">{errorMessage}</p>
          ) : lecture ? (
            <>
              <header className="detail-panel__title">
                <h2 id="lecture-detail-title">{lecture.Title || '-'}</h2>
                <p>{lecture.EnglishTitle || ''}</p>
              </header>
              <div className="detail-grid">
                <dl>
                  <dt>開講元</dt>
                  <dd>{lecture.Department || '-'}</dd>
                </dl>
                <dl>
                  <dt>担当教員</dt>
                  <dd>{formatTeachers(lecture.Teachers)}</dd>
                </dl>
                <dl>
                  <dt>授業形態</dt>
                  <dd>{formatLectureType(lecture.LectureType)}</dd>
                </dl>
                <dl>
                  <dt>曜日・時限</dt>
                  <dd>
                    {lecture.Timetables?.length ? (
                      <ul>
                        {lecture.Timetables.map((timetable) => (
                          <li key={`${timetable.LectureID}-${timetable.DayOfWeek}-${timetable.Period}`}>
                            {formatTimetable(timetable)}
                          </li>
                        ))}
                      </ul>
                    ) : (
                      '-'
                    )}
                  </dd>
                </dl>
                <dl>
                  <dt>開講時期</dt>
                  <dd>{lecture.Year ? `${lecture.Year}年度` : '-'}</dd>
                </dl>
                <dl>
                  <dt>科目コード</dt>
                  <dd>{lecture.Code || '-'}</dd>
                </dl>
                <dl>
                  <dt>単位数</dt>
                  <dd>{lecture.Credit || '-'}</dd>
                </dl>
                <dl>
                  <dt>対象学年</dt>
                  <dd>{formatLevelLabel(lecture.Level)}</dd>
                </dl>
                <dl>
                  <dt>使用言語</dt>
                  <dd>{lecture.Language || '-'}</dd>
                </dl>
                <dl>
                  <dt>URL</dt>
                  <dd>
                    {lecture.Url ? (
                      <a href={lecture.Url} target="_blank" rel="noreferrer">
                        シラバスを開く
                      </a>
                    ) : (
                      '-'
                    )}
                  </dd>
                </dl>
              </div>

              {splitByBreak(lecture.Abstract).length > 0 && (
                <section className="detail-section">
                  <h3>講義の概要とねらい</h3>
                  {splitByBreak(lecture.Abstract).map((line, index) => (
                    <p key={`abstract-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {lecture.LecturePlans?.length ? (
                <section className="detail-section">
                  <h3>授業計画・課題</h3>
                  <table>
                    <thead>
                      <tr>
                        <th>回</th>
                        <th>授業計画</th>
                        <th>課題</th>
                      </tr>
                    </thead>
                    <tbody>
                      {lecture.LecturePlans.map((plan, index) => (
                        <tr key={`${plan.Count}-${index}`}>
                          <td>{plan.Count || index + 1}</td>
                          <td>{plan.Plan || '-'}</td>
                          <td>{plan.Assignment || '-'}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </section>
              ) : null}

              {splitByBreak(lecture.Goal).length > 0 && (
                <section className="detail-section">
                  <h3>到達目標</h3>
                  {splitByBreak(lecture.Goal).map((line, index) => (
                    <p key={`goal-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {splitByBreak(lecture.Experience).length > 0 && (
                <section className="detail-section">
                  <h3>実務経験のある教員による授業</h3>
                  {splitByBreak(lecture.Experience).map((line, index) => (
                    <p key={`experience-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {lecture.Keywords?.length ? (
                <section className="detail-section">
                  <h3>キーワード</h3>
                  <ul className="detail-keywords">
                    {lecture.Keywords.map((keyword) => (
                      <li key={keyword}>{keyword}</li>
                    ))}
                  </ul>
                </section>
              ) : null}

              {splitByBreak(lecture.Textbook).length > 0 && (
                <section className="detail-section">
                  <h3>教科書</h3>
                  <ul>
                    {splitByBreak(lecture.Textbook).map((textbook, index) => (
                      <li key={`textbook-${index}`}>{textbook}</li>
                    ))}
                  </ul>
                </section>
              )}

              {splitByBreak(lecture.ReferenceBook).length > 0 && (
                <section className="detail-section">
                  <h3>参考書・講義資料等</h3>
                  <ul>
                    {splitByBreak(lecture.ReferenceBook).map((reference, index) => (
                      <li key={`reference-${index}`}>{reference}</li>
                    ))}
                  </ul>
                </section>
              )}

              {splitByBreak(lecture.Flow).length > 0 && (
                <section className="detail-section">
                  <h3>授業の進め方</h3>
                  {splitByBreak(lecture.Flow).map((line, index) => (
                    <p key={`flow-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {splitByBreak(lecture.OutOfClassWork).length > 0 && (
                <section className="detail-section">
                  <h3>授業時間外学修（予習・復習等）</h3>
                  {splitByBreak(lecture.OutOfClassWork).map((line, index) => (
                    <p key={`outofclass-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {splitByBreak(lecture.Assessment).length > 0 && (
                <section className="detail-section">
                  <h3>成績評価の基準及び方法</h3>
                  {splitByBreak(lecture.Assessment).map((line, index) => (
                    <p key={`assessment-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {lecture.RelatedCourses?.length ? (
                <section className="detail-section">
                  <h3>関連する科目</h3>
                  <ul>
                    {lecture.RelatedCourses.map((courseId) => (
                      <li key={`related-${courseId}`}>{courseId}</li>
                    ))}
                  </ul>
                </section>
              ) : null}

              {splitByBreak(lecture.Prerequisite).length > 0 && (
                <section className="detail-section">
                  <h3>履修の条件</h3>
                  {splitByBreak(lecture.Prerequisite).map((line, index) => (
                    <p key={`prerequisite-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {splitByBreak(lecture.Note).length > 0 && (
                <section className="detail-section">
                  <h3>その他</h3>
                  {splitByBreak(lecture.Note).map((line, index) => (
                    <p key={`note-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {splitByBreak(lecture.Contact).length > 0 && (
                <section className="detail-section">
                  <h3>連絡先</h3>
                  {splitByBreak(lecture.Contact).map((line, index) => (
                    <p key={`contact-${index}`}>{line}</p>
                  ))}
                </section>
              )}

              {splitByBreak(lecture.OfficeHours).length > 0 && (
                <section className="detail-section">
                  <h3>オフィスアワー</h3>
                  {splitByBreak(lecture.OfficeHours).map((line, index) => (
                    <p key={`office-${index}`}>{line}</p>
                  ))}
                </section>
              )}
            </>
          ) : (
            <p className="detail-loading">講義詳細が見つかりませんでした。</p>
          )}
        </div>
      </aside>
  );
};

export default CourseDetailPanel;
