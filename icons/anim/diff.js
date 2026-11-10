import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const plus = el.querySelectorAll('path[d="M12 3v14"], path[d="M5 10h14"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid });

  const minus = el.querySelector('path[d="M5 21h14"]');
  animate(minus, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
