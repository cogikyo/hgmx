import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const start = el.querySelector('path[d^="M12 7H7"]');
  animate(start, { translateX: [0, -1.5, 0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });

  const end = el.querySelector('path[d^="M12 17h5"]');
  animate(end, { translateX: [0, 1.5, -0.5, 0], translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
}
