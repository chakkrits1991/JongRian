import { useEffect, useState } from 'react'
import axios from 'axios'
import { Calendar, Clock, PlusCircle, CheckCircle, GraduationCap, BookmarkCheck } from 'lucide-react'

function App() {
  const [slots, setSlots] = useState([])
  const [startTime, setStartTime] = useState('')
  const [endTime, setEndTime] = useState('')

  const fetchSlots = () => {
    axios.get('http://localhost:8080/api/slots')
      .then(res => setSlots(res.data))
      .catch(err => console.error("Fetch error:", err))
  }

  useEffect(() => {
    fetchSlots()
  }, [])

  const handleSubmit = (e) => {
    e.preventDefault()
    axios.post('http://localhost:8080/api/slots', {
      start_time: new Date(startTime).toISOString(),
      end_time: new Date(endTime).toISOString()
    })
    .then(() => {
      alert("เพิ่มเวลาว่างเรียบร้อย!")
      setStartTime(''); setEndTime(''); fetchSlots()
    })
    .catch(err => alert("เกิดข้อผิดพลาด: " + err.response?.data?.error))
  }

  const handleBook = (id) => {
    if (!window.confirm("ยืนยันการจองคลาสนี้?")) return
    axios.patch(`http://localhost:8080/api/slots/${id}/book`)
      .then(() => {
        alert("จองสำเร็จ! เจอกันในคลาสครับ")
        fetchSlots()
      })
      .catch(err => alert(err.response?.data?.error || "จองไม่สำเร็จ"))
  }

  return (
    <div className="min-h-screen bg-slate-50 py-10 px-4 font-sans text-slate-900">
      <div className="max-w-5xl mx-auto">
        
        {/* Header Section */}
        <header className="flex items-center justify-between mb-12 bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
          <div className="flex items-center gap-4">
            <div className="bg-indigo-600 p-3 rounded-xl shadow-lg shadow-indigo-200">
              <GraduationCap className="text-white" size={32} />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight">JongRian System</h1>
              <p className="text-slate-500 text-sm">ระบบจองเวลาเรียนออนไลน์</p>
            </div>
          </div>
        </header>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-10">
          
          {/* ส่วนของ Tutor: เพิ่มเวลา */}
          <section className="lg:col-span-1">
            <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100 sticky top-10">
              <h2 className="text-lg font-bold mb-6 flex items-center gap-2 text-indigo-600">
                <PlusCircle size={20} /> สำหรับติวเตอร์
              </h2>
              <form onSubmit={handleSubmit} className="space-y-5">
                <div>
                  <label className="block text-xs font-bold text-slate-500 uppercase mb-2">เวลาเริ่มสอน</label>
                  <input 
                    type="datetime-local" 
                    className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
                    value={startTime} 
                    onChange={(e) => setStartTime(e.target.value)} 
                    required 
                  />
                </div>
                <div>
                  <label className="block text-xs font-bold text-slate-500 uppercase mb-2">เวลาเลิกสอน</label>
                  <input 
                    type="datetime-local" 
                    className="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all"
                    value={endTime} 
                    onChange={(e) => setEndTime(e.target.value)} 
                    required 
                  />
                </div>
                <button 
                  type="submit" 
                  className="w-full bg-indigo-600 text-white font-bold py-3 rounded-xl hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100 active:scale-95"
                >
                  เปิดตารางสอน
                </button>
              </form>
            </div>
          </section>

          {/* ส่วนของ Student: ดูและจองเวลา */}
          <section className="lg:col-span-2">
            <h2 className="text-xl font-bold mb-6 flex items-center gap-2 text-slate-800">
              <Calendar size={22} className="text-indigo-600" /> ตารางเวลาที่เปิดรับจอง
            </h2>
            
            <div className="grid gap-4">
              {slots.length > 0 ? slots.map((slot) => (
                <div key={slot.id} className="bg-white p-5 rounded-2xl border border-slate-100 shadow-sm flex items-center justify-between group hover:border-indigo-200 transition-all">
                  <div className="flex gap-4 items-center">
                    <div className={`p-3 rounded-full ${slot.status === 'available' ? 'bg-green-50' : 'bg-slate-100'}`}>
                      <Clock size={20} className={slot.status === 'available' ? 'text-green-600' : 'text-slate-400'} />
                    </div>
                    <div>
                      <div className="text-slate-800 font-semibold">
                        {new Date(slot.start_time).toLocaleString('th-TH', { 
                          dateStyle: 'medium', 
                          timeStyle: 'short' 
                        })}
                      </div>
                      <span className={`text-xs font-bold uppercase tracking-wider ${
                        slot.status === 'available' ? 'text-green-600' : 'text-slate-400'
                      }`}>
                        • {slot.status === 'available' ? 'เปิดจอง' : 'จองแล้ว'}
                      </span>
                    </div>
                  </div>

                  {slot.status === 'available' ? (
                    <button 
                      onClick={() => handleBook(slot.id)}
                      className="flex items-center gap-2 bg-emerald-500 text-white px-5 py-2.5 rounded-xl font-bold hover:bg-emerald-600 transition-all shadow-md shadow-emerald-100 active:scale-95"
                    >
                      <BookmarkCheck size={18} /> จองเลย
                    </button>
                  ) : (
                    <div className="text-slate-400 text-sm font-medium italic pr-4">ปิดรับจอง</div>
                  )}
                </div>
              )) : (
                <div className="text-center py-20 bg-white rounded-3xl border-2 border-dashed border-slate-200">
                  <p className="text-slate-400 font-medium text-lg">ยังไม่มีติวเตอร์เปิดตารางสอน</p>
                </div>
              )}
            </div>
          </section>

        </div>
      </div>
    </div>
  )
}

export default App