import { animate, stagger } from "../motion.js";

export default function (el) {
  const gusts = ["M9.8", "M17.5", "M12.8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(gusts, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
