import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M21 11"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pointer = el.querySelector('path[d^="M12.034"]');
  animate(pointer, { translateX: [0, -3, 1, 0], translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
