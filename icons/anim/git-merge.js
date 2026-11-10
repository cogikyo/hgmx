import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const line = el.querySelector("path");
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const target = el.querySelector('circle[cx="18"]');
  animate(target, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
