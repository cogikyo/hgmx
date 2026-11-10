import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const smoke = ["M6 21", "M11 21", "M16 21"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(smoke, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const alarm = el.querySelector("g");
  animate(alarm, { rotate: [0, -6, 4, 0] }, { duration: τmid, delay: τmin });
}
