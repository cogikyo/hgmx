import { animate, stagger } from "../motion.js";

export default function (el) {
  const tasks = ["M6 5h12", "M4 12h10", "M12 19h8"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(tasks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
