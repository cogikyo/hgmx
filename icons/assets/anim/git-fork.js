import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll("path");
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const left = el.querySelector('circle[cx="6"]');
  animate(left, { translateX: [6, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const right = el.querySelector('circle[cx="18"]');
  animate(right, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
