import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const strokes = ["M11.5 2", "M4 8L4 7", "M8 9H15", "M4 15H20"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(strokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.09) });
}
