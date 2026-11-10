import { animate, τmid } from "../motion.js";

export default function (el) {
  const stem = el.querySelector('line[y1="4"][y2="20"]');
  animate(stem, { rotate: [-20, 5, 0] }, { duration: τmid });
  const top = el.querySelector('line[y1="4"][y2="4"]');
  animate(top, { translateX: [-3, 1, 0] }, { duration: τmid });
  const bottom = el.querySelector('line[y1="20"]');
  animate(bottom, { translateX: [3, -1, 0] }, { duration: τmid });
}
