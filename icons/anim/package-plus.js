import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = ["M12 22V12", "M21 1", "M3.29 7", "m7.5 4.27"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M16 17h6"], path[d="M19 14v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
