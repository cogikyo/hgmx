import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M2.992"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const open = el.querySelector('path[d^="m10 9"]');
  animate(open, { translateX: [2, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const close = el.querySelector('path[d^="m14 15"]');
  animate(close, { translateX: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
