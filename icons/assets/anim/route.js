import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const start = el.querySelector('circle[cx="6"]');
  animate(start, { scale: [0, 1.25, 1] }, { duration: τmin });

  const path = el.querySelector("path");
  animate(path, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const end = el.querySelector('circle[cx="18"]');
  animate(end, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
