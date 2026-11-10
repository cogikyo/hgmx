import { animate, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d="M3 5h18"]');
  animate(top, { translateX: [-4, 1, 0] }, { duration: 0.3 });
  const middle = el.querySelector('path[d="M3 12h18"]');
  animate(middle, { translateX: [4, -1, 0] }, { duration: 0.3, delay: 0.06 });
  const last = el.querySelector('path[d="M7 19h10"]');
  animate(last, { scaleX: [0, 1.25, 1] }, { duration: 0.3, delay: τmin });
}
