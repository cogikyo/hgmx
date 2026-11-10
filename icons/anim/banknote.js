import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const note = el.querySelector("rect");
  animate(note, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const seal = el.querySelector("circle");
  animate(seal, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });

  const dots = el.querySelector("path");
  animate(dots, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
