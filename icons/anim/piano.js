import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const keys = ["M6", "M10", "M14", "M18"].map((x) => el.querySelector(`path[d="${x} 13v4"]`));
  animate(keys, { scaleY: [1, 1.5, 1] }, { duration: τmin, delay: stagger(0.06) });
}
