import { ChangeEvent, FormEvent, useState } from 'react';
import './App.css';
import { Scrape, SearchLectures, GetLectureDetails } from '../wailsjs/go/main/App';
import { domain } from '../wailsjs/go/models';
import {
    DEPARTMENT_OPTIONS,
    LEVEL_OPTIONS,
    LEVEL_LABELS,
    LECTURE_TYPE_LABELS,
    YEAR_OPTIONS,
    parseKeywordInput
} from './constants';

type SearchPanelProps = {
    title: string;
    onTitleChange: (value: string) => void;
    teacherName: string;
    onTeacherChange: (value: string) => void;
    keywords: string;
    onKeywordsChange: (value: string) => void;
    year: number | '';
    onYearChange: (value: number | '') => void;
    departments: string[];
    onDepartmentsChange: (value: string[]) => void;
    levels: number[];
    onLevelToggle: (value: number) => void;
    onClear: () => void;
    onSearch: () => Promise<void> | void;
    onScrape: () => Promise<void> | void;
    searching: boolean;
    scraping: boolean;
};

type ResultsTableProps = {
    results: domain.LectureSummary[];
    loading: boolean;
    errorMessage: string | null;
    statusMessage: string | null;
    selectedId: number | null;
    onSelect: (id: number) => void;
};

type LectureDetailProps = {
    lecture: domain.Lecture | null;
    loading: boolean;
    errorMessage: string | null;
};

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
    const semester = SEMESTER_LABELS[timetable.Semester?.toLowerCase?.() ?? ''] ?? timetable.Semester;
    const day = DAY_OF_WEEK_LABELS[timetable.DayOfWeek?.toLowerCase?.() ?? ''] ?? timetable.DayOfWeek;
    const period = timetable.Period ? `${timetable.Period}限` : '';
    const room = timetable.Room?.Name ? `@${timetable.Room.Name}` : '';
    return [semester, day ? `${day}曜` : '', period, room].filter(Boolean).join(' ');
};

const formatTeachers = (teachers: domain.Teacher[]): string => teachers.map((teacher) => teacher.Name).join(', ');

const formatLectureType = (lectureType: string): string =>
    LECTURE_TYPE_LABELS[lectureType?.toLowerCase?.() ?? ''] ?? (lectureType || '未設定');

const formatLevel = (level: number): string => LEVEL_LABELS[level] ?? `レベル${level}`;

const SearchPanel = ({
    title,
    onTitleChange,
    teacherName,
    onTeacherChange,
    keywords,
    onKeywordsChange,
    year,
    onYearChange,
    departments,
    onDepartmentsChange,
    levels,
    onLevelToggle,
    onClear,
    onSearch,
    onScrape,
    searching,
    scraping
}: SearchPanelProps) => {
    const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        void onSearch();
    };

    const handleDepartmentChange = (event: ChangeEvent<HTMLSelectElement>) => {
        const selected = Array.from(event.target.selectedOptions).map((option) => option.value);
        onDepartmentsChange(selected);
    };

    const handleYearChange = (event: ChangeEvent<HTMLSelectElement>) => {
        const value = event.target.value;
        onYearChange(value === '' ? '' : Number(value));
    };

    return (
        <form className="search-panel" onSubmit={handleSubmit}>
            <h2 className="panel-title">検索条件</h2>
            <div className="form-field">
                <label className="form-label" htmlFor="title">
                    講義名
                </label>
                <input
                    id="title"
                    type="text"
                    value={title}
                    onChange={(event) => onTitleChange(event.target.value)}
                    placeholder="例: データサイエンス"
                />
            </div>
            <div className="form-field">
                <label className="form-label" htmlFor="teacher">
                    担当教員
                </label>
                <input
                    id="teacher"
                    type="text"
                    value={teacherName}
                    onChange={(event) => onTeacherChange(event.target.value)}
                    placeholder="教員名で絞り込み"
                />
            </div>
            <div className="form-field">
                <label className="form-label" htmlFor="keywords">
                    キーワード
                </label>
                <textarea
                    id="keywords"
                    value={keywords}
                    onChange={(event) => onKeywordsChange(event.target.value)}
                    placeholder="キーワードをカンマまたは空白で区切って入力"
                    rows={2}
                />
            </div>
            <div className="form-field">
                <label className="form-label" htmlFor="year">
                    年度
                </label>
                <select id="year" value={year === '' ? '' : String(year)} onChange={handleYearChange}>
                    <option value="">指定なし</option>
                    {YEAR_OPTIONS.map((option) => (
                        <option value={option} key={option}>
                            {option}年度
                        </option>
                    ))}
                </select>
            </div>
            <div className="form-field">
                <label className="form-label" htmlFor="departments">
                    開講元
                </label>
                <select
                    id="departments"
                    multiple
                    value={departments}
                    onChange={handleDepartmentChange}
                    size={Math.min(10, DEPARTMENT_OPTIONS.length)}
                >
                    {DEPARTMENT_OPTIONS.map((department) => (
                        <option value={department} key={department}>
                            {department}
                        </option>
                    ))}
                </select>
                <p className="field-description">Ctrl / Cmd キーを押しながら複数選択できます。</p>
            </div>
            <fieldset className="form-field levels">
                <legend className="form-label">対象学年</legend>
                <div className="checkbox-list">
                    {LEVEL_OPTIONS.map((option) => {
                        const checked = levels.includes(option.value);
                        return (
                            <label key={option.value} className="checkbox-item">
                                <input
                                    type="checkbox"
                                    checked={checked}
                                    onChange={() => onLevelToggle(option.value)}
                                />
                                <span>{option.label}</span>
                            </label>
                        );
                    })}
                </div>
            </fieldset>
            <div className="button-group">
                <button type="submit" disabled={searching} className="primary">
                    {searching ? '検索中...' : '検索'}
                </button>
                <button type="button" onClick={onClear} className="ghost">
                    条件をクリア
                </button>
            </div>
            <div className="scrape-panel">
                <p className="field-description">最新のシラバスを取得する場合はスクレイピングを実行してください。</p>
                <button type="button" onClick={() => void onScrape()} disabled={scraping} className="secondary">
                    {scraping ? 'スクレイピング中...' : 'スクレイピングを実行'}
                </button>
            </div>
        </form>
    );
};

const ResultsTable = ({
    results,
    loading,
    errorMessage,
    statusMessage,
    selectedId,
    onSelect
}: ResultsTableProps) => {
    return (
        <section className="results-card">
            <div className="card-header">
                <h2>検索結果</h2>
                {statusMessage && <span className="status-message">{statusMessage}</span>}
                {errorMessage && <span className="error-message">{errorMessage}</span>}
            </div>
            {loading ? (
                <p className="placeholder">検索結果を読み込み中です...</p>
            ) : results.length === 0 ? (
                <p className="placeholder">検索条件を入力して講義を探してください。</p>
            ) : (
                <div className="table-container">
                    <table className="results-table">
                        <thead>
                            <tr>
                                <th scope="col">講義名</th>
                                <th scope="col">開講元</th>
                                <th scope="col">年度</th>
                                <th scope="col">担当教員</th>
                                <th scope="col">時間割</th>
                            </tr>
                        </thead>
                        <tbody>
                            {results.map((summary) => {
                                const isSelected = summary.ID === selectedId;
                                return (
                                    <tr key={summary.ID} className={isSelected ? 'selected' : undefined}>
                                        <td>
                                            <button
                                                type="button"
                                                className="link-button"
                                                onClick={() => onSelect(summary.ID)}
                                                aria-pressed={isSelected}
                                            >
                                                <span className="lecture-title">{summary.Title}</span>
                                                <span className="lecture-code">{summary.Code}</span>
                                            </button>
                                        </td>
                                        <td>{summary.Department || '-'}</td>
                                        <td>{summary.Year || '-'}</td>
                                        <td>{summary.Teachers?.length ? formatTeachers(summary.Teachers) : '-'}</td>
                                        <td>
                                            {summary.Timetables?.length
                                                ? summary.Timetables.map(formatTimetable).join(', ')
                                                : '-'}
                                        </td>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            )}
        </section>
    );
};

const LectureDetail = ({ lecture, loading, errorMessage }: LectureDetailProps) => {
    if (loading) {
        return (
            <section className="detail-card">
                <h2 className="card-title">講義詳細</h2>
                <p className="placeholder">詳細を読み込み中です...</p>
            </section>
        );
    }

    if (errorMessage) {
        return (
            <section className="detail-card">
                <h2 className="card-title">講義詳細</h2>
                <p className="error-message">{errorMessage}</p>
            </section>
        );
    }

    if (!lecture) {
        return (
            <section className="detail-card">
                <h2 className="card-title">講義詳細</h2>
                <p className="placeholder">講義を選択すると詳細が表示されます。</p>
            </section>
        );
    }

    const keywords = lecture.Keywords ?? [];

    return (
        <section className="detail-card">
            <h2 className="card-title">{lecture.Title}</h2>
            <p className="detail-lead">{lecture.EnglishTitle}</p>
            <div className="detail-grid">
                <dl>
                    <div>
                        <dt>科目コード</dt>
                        <dd>{lecture.Code || '-'}</dd>
                    </div>
                    <div>
                        <dt>開講元</dt>
                        <dd>{lecture.Department || '-'}</dd>
                    </div>
                    <div>
                        <dt>年度</dt>
                        <dd>{lecture.Year || '-'}</dd>
                    </div>
                    <div>
                        <dt>講義形態</dt>
                        <dd>{formatLectureType(lecture.LectureType)}</dd>
                    </div>
                    <div>
                        <dt>対象学年</dt>
                        <dd>{lecture.Level ? formatLevel(lecture.Level) : '-'}</dd>
                    </div>
                    <div>
                        <dt>単位数</dt>
                        <dd>{lecture.Credit || '-'}</dd>
                    </div>
                    <div>
                        <dt>使用言語</dt>
                        <dd>{lecture.Language || '-'}</dd>
                    </div>
                    <div>
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
                    </div>
                </dl>
                <dl>
                    <div>
                        <dt>担当教員</dt>
                        <dd>{lecture.Teachers?.length ? formatTeachers(lecture.Teachers) : '-'}</dd>
                    </div>
                    <div>
                        <dt>時間割</dt>
                        <dd>
                            {lecture.Timetables?.length ? (
                                <ul className="bullet-list">
                                    {lecture.Timetables.map((timetable, index) => (
                                        <li key={`${timetable.LectureID}-${index}`}>{formatTimetable(timetable)}</li>
                                    ))}
                                </ul>
                            ) : (
                                '-'
                            )}
                        </dd>
                    </div>
                    <div>
                        <dt>キーワード</dt>
                        <dd>
                            {keywords.length ? (
                                <ul className="keyword-list">
                                    {keywords.map((keyword) => (
                                        <li key={keyword}>{keyword}</li>
                                    ))}
                                </ul>
                            ) : (
                                '-'
                            )}
                        </dd>
                    </div>
                </dl>
            </div>

            <section className="detail-section">
                <h3>講義概要</h3>
                <p>{lecture.Abstract || '記載がありません。'}</p>
            </section>
            {lecture.Goal && (
                <section className="detail-section">
                    <h3>到達目標</h3>
                    <p>{lecture.Goal}</p>
                </section>
            )}
            {lecture.Flow && (
                <section className="detail-section">
                    <h3>授業の進め方</h3>
                    <p>{lecture.Flow}</p>
                </section>
            )}
            {lecture.Assessment && (
                <section className="detail-section">
                    <h3>成績評価</h3>
                    <p>{lecture.Assessment}</p>
                </section>
            )}
            {lecture.LecturePlans?.length ? (
                <section className="detail-section">
                    <h3>講義計画</h3>
                    <table className="plans-table">
                        <thead>
                            <tr>
                                <th scope="col">回</th>
                                <th scope="col">内容</th>
                                <th scope="col">課題</th>
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
        </section>
    );
};

function App() {
    const [title, setTitle] = useState('');
    const [teacherName, setTeacherName] = useState('');
    const [keywords, setKeywords] = useState('');
    const [year, setYear] = useState<number | ''>('');
    const [departments, setDepartments] = useState<string[]>([]);
    const [levels, setLevels] = useState<number[]>([]);

    const [lectures, setLectures] = useState<domain.LectureSummary[]>([]);
    const [selectedLectureId, setSelectedLectureId] = useState<number | null>(null);
    const [selectedLecture, setSelectedLecture] = useState<domain.Lecture | null>(null);

    const [searching, setSearching] = useState(false);
    const [detailLoading, setDetailLoading] = useState(false);
    const [scraping, setScraping] = useState(false);

    const [searchError, setSearchError] = useState<string | null>(null);
    const [detailError, setDetailError] = useState<string | null>(null);
    const [statusMessage, setStatusMessage] = useState<string | null>(null);

    const handleLevelToggle = (value: number) => {
        setLevels((prev) => (prev.includes(value) ? prev.filter((item) => item !== value) : [...prev, value]));
    };

    const buildSearchQuery = (): domain.SearchQuery => {
        const keywordList = parseKeywordInput(keywords);
        const levelValues = levels.slice();
        const selectedYear = typeof year === 'number' ? year : 0;
        return domain.SearchQuery.createFrom({
            Title: title.trim(),
            Keywords: keywordList,
            Departments: departments,
            Year: selectedYear,
            TeacherName: teacherName.trim(),
            TimeTables: [],
            Levels: levelValues
        });
    };

    const handleSearch = async () => {
        setSearching(true);
        setSearchError(null);
        setStatusMessage(null);
        try {
            const query = buildSearchQuery();
            const result = await SearchLectures(query);
            setLectures(result);
            setSelectedLectureId(null);
            setSelectedLecture(null);
            setStatusMessage(result.length ? `${result.length} 件の講義が見つかりました。` : '講義は見つかりませんでした。');
        } catch (error) {
            console.error(error);
            setSearchError('検索に失敗しました。時間をおいて再度お試しください。');
        } finally {
            setSearching(false);
        }
    };

    const handleClear = () => {
        setTitle('');
        setTeacherName('');
        setKeywords('');
    setYear('');
        setDepartments([]);
        setLevels([]);
        setLectures([]);
        setSelectedLectureId(null);
        setSelectedLecture(null);
        setStatusMessage(null);
        setSearchError(null);
        setDetailError(null);
    };

    const loadLectureDetail = async (lectureId: number) => {
        setSelectedLectureId(lectureId);
        setDetailLoading(true);
        setDetailError(null);
        try {
            const lecture = await GetLectureDetails(lectureId);
            setSelectedLecture(lecture ?? null);
        } catch (error) {
            console.error(error);
            setDetailError('講義詳細の取得に失敗しました。');
        } finally {
            setDetailLoading(false);
        }
    };

    const handleScrape = async () => {
        setScraping(true);
        setStatusMessage(null);
        setSearchError(null);
        try {
            await Scrape();
            setStatusMessage('最新のシラバス情報を取得しました。必要に応じて再検索してください。');
        } catch (error) {
            console.error(error);
            setSearchError('スクレイピングに失敗しました。');
        } finally {
            setScraping(false);
        }
    };

    return (
        <div className="app">
            <header className="app-header">
                <div>
                    <h1>Science Tokyo シラバス検索</h1>
                    <p className="header-subtitle">講義の検索や詳細の確認、最新データの取得を行えます。</p>
                </div>
            </header>
            <main className="app-main">
                <SearchPanel
                    title={title}
                    onTitleChange={setTitle}
                    teacherName={teacherName}
                    onTeacherChange={setTeacherName}
                    keywords={keywords}
                    onKeywordsChange={setKeywords}
                    year={year}
                    onYearChange={setYear}
                    departments={departments}
                    onDepartmentsChange={setDepartments}
                    levels={levels}
                    onLevelToggle={handleLevelToggle}
                    onClear={handleClear}
                    onSearch={handleSearch}
                    onScrape={handleScrape}
                    searching={searching}
                    scraping={scraping}
                />
                <div className="results-wrapper">
                    <ResultsTable
                        results={lectures}
                        loading={searching}
                        errorMessage={searchError}
                        statusMessage={statusMessage}
                        selectedId={selectedLectureId}
                        onSelect={loadLectureDetail}
                    />
                    <LectureDetail lecture={selectedLecture} loading={detailLoading} errorMessage={detailError} />
                </div>
            </main>
        </div>
    );
}

export default App;
