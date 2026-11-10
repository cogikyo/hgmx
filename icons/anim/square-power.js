import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stem = el.querySelector('path[d="M12 7v4"]');
  animate(stem, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
