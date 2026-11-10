import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = ["M12 22V12", "M21 10.342", "M3.29 7", "m7.5 4.27"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mark = el.querySelectorAll('path[d^="M17 15.27"], path[d^="M19 21.99"]');
  animate(mark, { rotateY: [0, 360], scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });
}
