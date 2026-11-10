import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const parens = el.querySelectorAll("path");
  animate(parens, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const x = el.querySelectorAll("line");
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
