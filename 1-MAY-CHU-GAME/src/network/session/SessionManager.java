// /*    */ package network.session;
// /*    */ 
// /*    */ import java.util.ArrayList;
// /*    */ import java.util.List;
// /*    */ 
// /*    */ 
// /*    */ 
// /*    */ 
// /*    */ 
// /*    */ public class SessionManager
// /*    */ {
// /*    */   private static SessionManager i;
// /*    */   private List<Session> sessions;
// /*    */   
// /*    */   public static SessionManager gI() {
// /* 16 */     if (i == null) {
// /* 17 */       i = new SessionManager();
// /*    */     }
// /* 19 */     return i;
// /*    */   }
// /*    */ 
// /*    */ 
// /*    */   
// /*    */   public SessionManager() {
// /* 25 */     this.sessions = new ArrayList<>();
// /*    */   }
// /*    */ }


// /* Location:              C:\Users\VoHoangKiet\Downloads\TEA_V5\lib\GirlkunNetwork.jar!\com\girlkun\network\session\SessionManager.class
//  * Java compiler version: 8 (52.0)
//  * JD-Core Version:       1.1.3
//  */
package network.session;

import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;

/**
 * [FIX-5] volatile + double-checked locking
 *         Class gốc không thread-safe → 2 thread vào gI() cùng lúc tạo ra
 *         2 instance khác nhau → mỗi thread thấy một list sessions riêng
 *         → Sweep #2 trong login() không tìm thấy duplicate → không kick
 *         → 2 session cùng username tồn tại song song → crash chéo nhau
 *
 * [FIX-6] CopyOnWriteArrayList thay ArrayList
 *         ArrayList không thread-safe → crash khi add/remove/iterate đồng thời
 *
 * [FIX-7] Constructor private
 *         Bản gốc constructor public → có thể new SessionManager() từ nơi khác
 *         → tạo ra instance độc lập, sessions trong đó không ai biết
 *
 * [NEW]   Thêm addSession/removeSession/isConnectedSession
 *         Class gốc hoàn toàn trống, không có method nào → sessions không được
 *         quản lý → tích lũy session chết theo thời gian → crash sau vài tiếng
 */
public class SessionManager {

    // [FIX-5] volatile: đảm bảo mọi thread thấy giá trị mới nhất sau khi ghi
    private static volatile SessionManager i;

    // [FIX-6] CopyOnWriteArrayList: an toàn khi iterate + modify đồng thời
    private final List<Session> sessions;

    // [FIX-7] private constructor
    private SessionManager() {
        this.sessions = new CopyOnWriteArrayList<>();
    }

    // [FIX-5] Double-checked locking: chỉ synchronized khi thực sự cần tạo instance
    public static SessionManager gI() {
        if (i == null) {
            synchronized (SessionManager.class) {
                if (i == null) {
                    i = new SessionManager();
                }
            }
        }
        return i;
    }

    // [NEW] Gọi khi client kết nối vào — thêm vào list để quản lý
    public void addSession(Session session) {
        if (session != null && !sessions.contains(session)) {
            sessions.add(session);
        }
    }

    // [NEW] Gọi khi client disconnect/bị kick — dọn khỏi list ngay
    // Đây là bước quan trọng nhất: nếu không remove thì list tích lũy
    // session chết → Sweep #2 trong login() duyệt qua → cast → NPE
    public void removeSession(Session session) {
        sessions.remove(session);
    }

    // [NEW] Lấy danh sách — dùng trong Sweep #2 của login()
    // Trả về CopyOnWriteArrayList nên caller có thể iterate an toàn
    public List<Session> getSessions() {
        return sessions;
    }

    // [NEW] Tiện ích: số session đang sống
    public int getOnlineCount() {
        return sessions.size();
    }

    // [NEW] Tiện ích: dọn toàn bộ — dùng khi restart/maintenance
    public void removeAll() {
        sessions.clear();
    }
}