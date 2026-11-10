import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const columns = el.querySelectorAll('path[d$="v18"]');
  animate(columns, { translateX: [2, -1, 0] }, { duration: τmid });

  const rows = el.querySelectorAll('path[d$="h18"]');
  animate(rows, { translateY: [-2, 1, 0] }, { duration: τmid, delay: τmin });
}
