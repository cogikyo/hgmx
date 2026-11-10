import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const horns = el.querySelector("path");
  animate(horns, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bull = el.querySelector("g");
  animate(bull, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
