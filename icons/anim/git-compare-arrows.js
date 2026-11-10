import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="m15 9"]');
  animate(left, { translateX: [3, 0], opacity: [0, 1] }, { duration: τmid });

  const right = el.querySelector('path[d^="m9 15"]');
  animate(right, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
