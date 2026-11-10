import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const top = el.querySelector('path[d="M7 8h10"]');
  animate(top, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });

  const bottom = el.querySelector('path[d="M7 16h10"]');
  animate(bottom, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
