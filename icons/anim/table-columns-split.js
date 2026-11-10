import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cut = ["M14 2v2", "M14 8v2", "M14 14v2", "M14 20v2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(cut, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const left = el.querySelectorAll('path[d^="M2 "], path[d="M5 3v18"]');
  animate(left, { translateX: [0, -1.5, 0] }, { duration: τmid, delay: τmin });

  const right = el.querySelectorAll('path[d^="M22 "]');
  animate(right, { translateX: [0, 1.5, 0] }, { duration: τmid, delay: τmin });
}
