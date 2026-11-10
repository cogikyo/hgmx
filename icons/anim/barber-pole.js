import { animate, stagger } from "../motion.js";

export default function (el) {
  const stripes = ["m8 20.5", "m8 16", "m8 11.5"].map((d) => el.querySelector(`path[d^="${d} "]`));
  animate(stripes, { pathLength: [0, 1], opacity: [0, 1], translateY: [1.5, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
