import { animate, stagger } from "../motion.js";

export default function (el) {
  const keys = ["M6", "M10", "M14", "M18"].map((x) => el.querySelector(`path[d="${x} 13v4"]`));
  animate(keys, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
