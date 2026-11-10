import { animate, stagger } from "../motion.js";

export default function (el) {
  const events = ["M9.414 5.414", "M9.414 13.414", "M9.414 21.414"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(events, { translateX: [0, 2, -1, 0] }, { duration: 0.3, delay: stagger(0.09) });
}
