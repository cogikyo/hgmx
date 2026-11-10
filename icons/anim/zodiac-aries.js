import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const stem = el.querySelector('path[d^="M7 12"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const horn = el.querySelector('path[d^="M12 7.5"]');
  animate(horn, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
